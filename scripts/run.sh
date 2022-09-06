#!/usr/bin/env bash

set -eo pipefail

__DIR="$(cd "$(dirname "$0")" && pwd)"

function usage() {
    echo "Usage:"
    echo "  $0 [flags]"
    printf "\n"
    echo "Flags:"
    printf "  %s, --help\t\tPrints usage\n" "-h"
    printf "  %s, --user string\tThe BOSH ssh user [default: vcap]\n" "-u"
    printf "  %s, --private-key-path string\tThe path to private key when using the bbr user\n" "-i"
    printf "\n"
}

USER=vcap

while [ "$1" != "" ]; do
    param=$(echo "$1" | awk -F= '{print $1}')
    value=$(echo "$1" | awk -F= '{print $2}')
    case $param in
      -h | --help)
        usage
        exit
        ;;
      -u | --user)
        USER=$value
        ;;
      -i | --private-key-path)
        PRIVATE_KEY_PATH=$value
        ;;
      help)
        usage
        exit
        ;;
      *)
        echo ""
        echo "Invalid option: [$param]"
        echo ""
        usage
        exit 1
        ;;
    esac
    shift
done

if [[ -z "${DEPLOYMENT}" ]]; then
    echo -n "Enter cf deployment name: "
    read -r DEPLOYMENT
fi

if [[ -z "${BOSH_VM_IP}" ]]; then
    echo -n "Enter the IP for bosh director: "
    read -r BOSH_VM_IP
fi

if [[ "${USER}" == "vcap" && -z "${VCAP_PASSWORD}" ]]; then
    echo -n "Enter password for vcap: "
    read -rs VCAP_PASSWORD
fi

if [[ "${USER}" == "bbr" && -z "${PRIVATE_KEY_PATH}" ]]; then
    echo -n "Enter path to bbr private key: "
    read -r PRIVATE_KEY_PATH
fi

script_dir=$(mktemp -d)

cat  > "$script_dir/install.sh" <<EOF
#!/bin/bash

set -eo pipefail

export http_proxy=$HTTP_PROXY
export https_proxy=$HTTPS_PROXY
export no_proxy=$NO_PROXY

BIN_PATH="/home/$USER/bin"
EOF
cat "$__DIR"/install.sh >> "$script_dir/install.sh"

cat > "$script_dir/pdm.sh" <<EOF
#!/bin/bash

/home/$USER/bin/pdm -n $DEPLOYMENT
EOF

if [[ "$USER" == "bbr" ]]; then
    ssh-keygen -f "$HOME/.ssh/known_hosts" -R "$BOSH_VM_IP"
    chmod 600 "$PRIVATE_KEY_PATH"
    if ssh -o StrictHostKeyChecking=no -i "$PRIVATE_KEY_PATH" "bbr@$BOSH_VM_IP" -q "bash -s " <  "$script_dir/install.sh" > /dev/null 2>&1; then
        ssh -o StrictHostKeyChecking=no -i "$PRIVATE_KEY_PATH" "bbr@$BOSH_VM_IP" -q "bash -s " <  "$script_dir/pdm.sh"
    else
        echo "failed to execute $script_dir/install.sh"
    fi
else
    if sshpass -p "$VCAP_PASSWORD" ssh "vcap@$BOSH_VM_IP" -q "bash -s " <  "$script_dir/install.sh" > /dev/null 2>&1; then
        sshpass -p "$VCAP_PASSWORD" ssh "vcap@$BOSH_VM_IP" -q "bash -s " <  "$script_dir/pdm.sh"
    else
        echo "failed to execute $script_dir/install.sh"
    fi
fi
