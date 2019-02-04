package router

import (
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	RegTestClient *RegTest
	LiquidClient  *Liquid
}

// New creates a new router
func New(client *RegTest, liquidClient *Liquid) *Router {
	router := mux.NewRouter().StrictSlash(true)
	r := &Router{router, client, liquidClient}

	r.HandleFunc("/send", r.RegTestClient.SendTo).Methods("POST")
	r.HandleFunc("/broadcast", r.RegTestClient.Broadcast).Methods("POST")
	r.HandleFunc("/utxos/{address}", r.RegTestClient.GetUtxos).Methods("GET")
	r.HandleFunc("/fees", r.RegTestClient.EstimateFees).Methods("GET")
	r.HandleFunc("/ping", r.RegTestClient.Ping).Methods("GET")
	r.HandleFunc("/txs/{hash}", r.RegTestClient.GetTx).Methods("GET")

	//Liquid
	r.HandleFunc("/liquid/ping", r.LiquidClient.Ping).Methods("GET")
	r.HandleFunc("/liquid/block", r.LiquidClient.BlockCount).Methods("GET")
	//NOTICE: This uses Liquid instance wallet's db
	r.HandleFunc("/liquid/issue", r.LiquidClient.ListIssuances).Methods("GET")
	r.HandleFunc("/liquid/issue", r.LiquidClient.IssueFractionalAsset).Methods("POST")
	//r.HandleFunc("/liquid/send", r.LiquidClient.SendAsset).Methods("POST")

	return r
}
