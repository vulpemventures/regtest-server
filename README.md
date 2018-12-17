# Regtest server

This repo contains a Go implementation of a web server exposing the following endpoints to interact with the underline cluster of bitcoin daemons in regtest mode:

* `send/:address` sends funds to the target address and returns the transaction hash
* `broadcast/:tx` publishes the passed raw tx to the regtest network and returns the transaction hash
* `utxos/:address` returns all **unspent** transaction output related to the target address [WIP]

It uses [bitcoin-testnet-box](https://github.com/vulpemventures/bitcoin-testnet-box/) fork.

You can use [Docker](https://docker.io) to start a cluster of bitcoin daemons or a linux installer bash script [here](https://github.com/vulpemventures/regtest-server/blob/master/scripts/install-bitcoin)

## Clone

Bitcoin cluster

```sh
git clone https://github.com/vulpemventures/bitcoin-testnet-box.git
```

Regtest server

```sh
git clone https://github.com/vulpemventures/regtest-server.git
```

## Run bitcoin cluster

Enter `bitcoin-testnet-box` folder and run

```sh
make start
```


## Build regtest server

Enter `regtest-server` folder

```sh
./scripts/buildlinux
```


Start server in another tab

```sh
# run server at http://localhost:8000/
./build/regtest-server-linux-amd64
# or specify url
ADDRESS=192.168.0.20 PORT=8001 ./build/regtest-server-linux-amd64
```

## Stop

In the `bitcoin-testnet-box` folder run

```sh
make stop
make clean
exit
```
