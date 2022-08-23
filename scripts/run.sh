#!/usr/bin/env bash

set -eo pipefail

__DIR="$(cd "$(dirname "$0")" && pwd)"

if [[ -z "${DEPLOYMENT}" ]]; then
    echo -n "Enter cf deployment name: "
    read -r DEPLOYMENT
fi

if [[ -z "${BOSH_VM_IP}" ]]; then
    echo -n "Enter the IP for bosh director: "
    read -r BOSH_VM_IP
fi

if [[ -z "${VCAP_PASSWORD}" ]]; then
    echo -n "Enter password for vcap: "
    read -rs VCAP_PASSWORD
fi

script_dir=$(mktemp -d)

cat  > "$script_dir/install.sh" <<EOF
#!/bin/bash

set -eo pipefail

export http_proxy=$HTTP_PROXY
export https_proxy=$HTTPS_PROXY
export no_proxy=$NO_PROXY

EOF
cat $__DIR/install.sh >> "$script_dir/install.sh"

cat > "$script_dir/pdm.sh" <<EOF
#!/bin/bash

/home/vcap/bin/pdm -n $DEPLOYMENT
EOF

if sshpass -p $VCAP_PASSWORD ssh "vcap@$BOSH_VM_IP" -q "bash -s " <  "$script_dir/install.sh" > /dev/null 2>&1; then
    sshpass -p $VCAP_PASSWORD ssh "vcap@$BOSH_VM_IP" -q "bash -s " <  "$script_dir/pdm.sh"
else
    echo "failed to execute install.sh"
fi

