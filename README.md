# bosh-persistent-disk-modifier

Run this tool to update the disk CID mappings in the bosh database after an HCX
migration of CF or other bosh deployment.

## Setup

1. Export the following variables

    ```sh
    export BOSH_VM_IP=10.0.0.21
    export VCAP_PASSWORD='bosh vm vcap password'
    export DEPLOYMENT=cf-02614dc53e91b381e7bd
    ```

1. If necessary, export http proxy variables

    ```sh
    export HTTP_PROXY=http://some.proxy.local
    export HTTPS_PROXY=http://some.proxy.local
    export NO_PROXY=comma-delimitted-excluded-ips-domains-from-proxy
    ```

## Run

    ```sh
    ./scripts/run.sh
    ```sh
