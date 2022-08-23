#!/usr/bin/env bash

set -eo pipefail

__DIR="$(cd "$(dirname "$0")" && pwd)"

if [[ -z "${DEPLOYMENT}" ]]; then
    echo -n "Enter cf deployment name: "
    read -r DEPLOYMENT
fi

if [[ -z "${VCAP_PASSWORD}" ]]; then
    echo -n "Enter password for vcap: "
    read -rs VCAP_PASSWORD
fi

script_dir=$(mktemp -d)
cat > "$script_dir/pdm.sh" <<EOF
#!/bin/bash

/home/vcap/bin/pdm -n $DEPLOYMENT
EOF

if sshpass -p $VCAP_PASSWORD ssh vcap@10.176.168.22 -q "bash -s " <  "$__DIR/install.sh" > /dev/null 2>&1; then
    sshpass -p $VCAP_PASSWORD ssh vcap@10.176.168.22 -q "bash -s " <  "$script_dir/pdm.sh"
fi

