package router

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var (
	LiquidHost     = flag.String("LIQUID_HOST", "localhost", "Host of Liquid instance")
	LiquidPort     = flag.Int("LIQUID_PORT", 18884, "Port of Liquid 18884")
	LiquidUser     = flag.String("LIQUID_RPC_USER", "user1", "Liquid RPC user")
	LiquidPassword = flag.String("LIQUID_RPC_PASSWORD", "password1", "Liquid RPC password")
	SSL            = false
)

// Liquid handles communication with both the regtest daemon and the local database
// @param Client <*Client>: handles remote API calls to the regtest daemon
// @param DB <*Database>: BBolt database is used to keep track of the utxos.
// 	This handles read/write/list/delete operations to the db
type Liquid struct {
	//DB *Database
	client *rpcClient
}

// New configures and creates a Client instance
func (lqd *Liquid) New() error {

	jsonRPCClient, err := newRpcClient(*LiquidHost, *LiquidPort, *LiquidUser, *LiquidPassword, SSL, 6)
	if err != nil {
		log.Fatalln(err)
	}

	lqd.client = jsonRPCClient

	// Setup db
	/* 	db := &Database{}
	   	err := db.New()
	   	if err != nil {
	   		return err
	   	}

	   	lqd.DB = db
	*/
	return nil

}

//Ping is a function to ping
func (lqd *Liquid) Ping(w http.ResponseWriter, req *http.Request) {
	r, err := lqd.client.call("getblockchaininfo", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}

// BlockCount returns he block count
func (lqd *Liquid) BlockCount(w http.ResponseWriter, req *http.Request) {
	r, err := lqd.client.call("getblockcount", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}

// IssueFractionalAsset mints to the given address 100 000 00 fractional untis that represents
// the single offline asset held in custody by a trusted authority
func (lqd *Liquid) IssueFractionalAsset(w http.ResponseWriter, req *http.Request) {

	resp, err := lqd.client.call("issueasset", []interface{}{1, 0, false})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Mine one block
	_, err = lqd.client.call("generate", []interface{}{1})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

//SendAsset send to the given address, the given asset type and amount
/* func (lqd *Liquid) SendAsset(w http.ResponseWriter, req *http.Request) {
	body := getRequestBody(req.Body)
	address := body["address"]
	asset := body["asset"]
	amount := body["amount"]
	subtractfeefromamount := true
	var comment string
	var commentTo string

	r, err := lqd.client.call("sendtoaddress", []interface{}{address, amount, comment, commentTo, subtractfeefromamount, asset})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Mine one block
	_, err = lqd.client.call("generate", []interface{}{1})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
} */

//ListIssuances list all asset being issued by the liquid instance
func (lqd *Liquid) ListIssuances(w http.ResponseWriter, req *http.Request) {
	//listissuances
	r, err := lqd.client.call("listissuances", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
