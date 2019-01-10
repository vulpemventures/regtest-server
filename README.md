# Regtest server

This repo contains a Go implementation of a web server exposing the following endpoints to interact with the underline cluster of bitcoin daemons in regtest mode:

* `POST send/ {"address": "receving_address"}` sends funds to the receving address and returns the transaction hash
* `POST broadcast/ {"tx": "signed_transaction_hex"}` publishes the passed raw tx to the regtest network and returns the transaction hash
* `GET utxos/:address` returns all **unspent** transaction outputs related to the target address
* `GET txs/:txhash` returns detailed info about the transaction with the corresponding hash
* `GET fees/` returns minimum mandatory fee amount to include in the transaction to be mined
* `GET ping/` returns info about the service, used as health check endpoint

It uses [bitcoin-testnet-box](https://github.com/vulpemventures/bitcoin-testnet-box/) fork.

You can use [Docker](https://docker.io) to start a cluster of bitcoin daemons or a linux installer bash script [here](https://github.com/vulpemventures/regtest-server/blob/master/scripts/install-bitcoin)

## Clone

Bitcoin cluster

Clone this repo if you don't want to use Docker

```sh
git clone https://github.com/vulpemventures/bitcoin-testnet-box.git
```

Regtest server

```sh
git clone https://github.com/vulpemventures/regtest-server.git
```

## Run cluster without Docker

Enter `bitcoin-testnet-box` folder and run

```sh
make start
make generate BLOCKS=200
```

## Run cluster with Docker

Enter `regtest-server` folder and run

```sh
./scripts/run_server.sh

# inside docker shell
make start
make generate BLOCKS=200
```

## Build regtest server (linux)

Enter `regtest-server` folder

```sh
./scripts/buildlinux
```

## Build regtest server (mac)

Enter `regtest-server` folder

```sh
./scripts/build darwin amd64
```

Start server in another tab

```sh
# run server at http://localhost:8000/
./build/regtest-server-linux-amd64
# or specify url
ADDRESS=192.168.0.20 PORT=8001 ./build/regtest-server-linux-amd64
```

## Stop

In the `bitcoin-testnet-box` folder or in the docker shell run

```sh
make stop
make clean
exit
```
