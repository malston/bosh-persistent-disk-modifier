#!/bin/bash

set -eo pipefail

VERSION="latest"
OS=linux
ARCH=amd64
BIN_PATH="/usr/local/bin"

function usage() {
    echo "Usage:"
    echo "  $0 [flags]"
    printf "\n"
    echo "Flags:"
    printf "  %s, --help\t\tPrints usage\n" "-h"
    printf "  %s, --version string\tVersion [default: latest]\n" "-v"
    printf "  %s, --os string\tOperating system [default: linux]\n" "-o"
    printf "  %s, --arch string\tArchitecture [default: amd64]\n" "-a"
    printf "  %s, --path string\tInstall path [default: /usr/local/bin]\n" "-p"
    printf "\n"
    echo "Examples:"
    printf "  %s --version=0.0.1 " "$0"
    printf "\n"
}

function install_pdm() {
    version="${1}"
    os="${2}"
    arch="${3}"
    bin_path="${4}"
    file="pdm.tgz"
    trap '{ rm -f "$file" ; exit 0; }' EXIT
    if [[ $version = "latest" ]]; then
        curl -fsSL -o $file "https://github.com/malston/bosh-persistent-disk-modifier/releases/latest/download/pdm-$os-$arch.tgz"
    elif [[ $version == v* ]]; then
        curl -fsSL -o $file "https://github.com/malston/bosh-persistent-disk-modifier/releases/download/$version/pdm-$os-$arch.tgz"
    else
        curl -fsSL -o $file "https://github.com/malston/bosh-persistent-disk-modifier/releases/download/v$version/pdm-$os-$arch.tgz"
    fi
    mkdir -p "$bin_path"
    tar -xvf $file -C "$bin_path"
    chmod +x "$bin_path/pdm"
}

while [ "$1" != "" ]; do
    param=$(echo "$1" | awk -F= '{print $1}')
    value=$(echo "$1" | awk -F= '{print $2}')
    case $param in
      -h | --help)
        usage
        exit
        ;;
      -o | --os)
        OS=$value
        ;;
      -a | --arch)
        ARCH=$value
        ;;
      -v | --version)
        VERSION=$value
        ;;
      -p | --path)
        BIN_PATH=$value
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

install_pdm "$VERSION" "$OS" "$ARCH" "$BIN_PATH"
