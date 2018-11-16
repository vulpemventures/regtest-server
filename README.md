# Regtest server

This repo contains a Go implementation of a web server exposing the following endpoints to interact with the underline cluster of bitcoin daemons in regtest mode:

* `send/:address` sends funds to the target address and returns the transaction hash
* `broadcast/:tx` publishes the passed raw tx to the regtest network and returns the transaction hash
* `utxos/:address` returns all **unspent** transaction output related to the target address [WIP]

It uses [bitcoin-testnet-box](https://github.com/freewil/bitcoin-testnet-box/) with [Docker](https://docker.io) to start the daemons.

## Install

Clone repo

```sh
git clone git@github.com:vulpemventures/regtest-server.git
cd regtest-server
```

Install dependencies and run

```sh
go get -d
./scripts/install_docker.sh
```

## Run regtest

Start daemon

```sh
./scripts/run_server.sh
make start
```

Start server in another tab

```sh
go build
# run server at http://localhost:8000/
./regtest-server
# or specify url
ADDRESS=192.168.0.20 PORT=8001 ./regtest-server
```

## Stop regtest

In the regtest tab run

```sh
make stop
make clean
exit
```

## Endpoint responses