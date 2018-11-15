package router

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gorilla/mux"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

type RegTest struct {
	Client        *rpcclient.Client
	FaucetAccount string
}

// New configures and creates a Client instance
func (r *RegTest) New() error {
	connConfig := &rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "localhost:19001",
		User:         "admin1",
		Pass:         "123",
	}

	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		return err
	}
	addr, err := client.GetAccountAddress("wallet.dat")
	if err != nil {
		return err
	}

	r.Client = client
	r.FaucetAccount = addr.String()

	mine(r, 200)

	return nil
}

// Shutdown disconnect from rpc server and stop all goroutines
func (r *RegTest) Shutdown() {
	r.Client.Shutdown()
}

// SendTo sends 1 btc to the given address from the faucet account
func (r *RegTest) SendTo(w http.ResponseWriter, req *http.Request) {
	address := mux.Vars(req)["address"]
	txHash, err := sendTo(r, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO save tx in db

	resp := fmt.Sprintf(`{"txHash": %s}`, txHash.String())
	json.NewEncoder(w).Encode(resp)
}

// Broadcast publishes a raw transaction to the network
func (r *RegTest) Broadcast(w http.ResponseWriter, req *http.Request) {
	tx := mux.Vars(req)["tx"]
	rawTx, err := hex.DecodeString(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	txHash, err := broadcast(r, rawTx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: save tx in db && update spent utxos

	resp := fmt.Sprintf(`{"txHash": %s}`, txHash.String())
	json.NewEncoder(w).Encode(resp)
}

func mine(r *RegTest, num int) ([]*chainhash.Hash, error) {
	return r.Client.Generate(uint32(num))
}

func sendTo(r *RegTest, address string) (*chainhash.Hash, error) {
	receiver, err := btcutil.DecodeAddress(address, &chaincfg.RegressionNetParams)
	if err != nil {
		return nil, err
	}
	txHash, err := r.Client.SendToAddress(receiver, btcutil.Amount(100000000))
	if err != nil {
		return nil, err
	}

	_, err = mine(r, 1)
	if err != nil {
		return nil, err
	}

	return txHash, nil
}

func broadcast(r *RegTest, tx []byte) (*chainhash.Hash, error) {
	rawTx := &wire.MsgTx{}
	err := rawTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		return nil, err
	}

	txHash, err := r.Client.SendRawTransaction(rawTx, true)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = mine(r, 1)
	if err != nil {
		return nil, err
	}

	return txHash, nil
}
