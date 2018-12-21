package router

import (
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	RegTestClient *RegTest
}

// New creates a new router
func New(client *RegTest) *Router {
	router := mux.NewRouter().StrictSlash(true)

	r := &Router{router, client}

	r.HandleFunc("/send", r.RegTestClient.SendTo).Methods("POST")
	r.HandleFunc("/broadcast", r.RegTestClient.Broadcast).Methods("POST")
	r.HandleFunc("/utxos/{address}", r.RegTestClient.GetUtxos).Methods("GET")
	r.HandleFunc("/fees", r.RegTestClient.EstimateFees).Methods("GET")
	r.HandleFunc("/ping", r.RegTestClient.Ping).Methods("GET")
	r.HandleFunc("/txs/{hash}", r.RegTestClient.GetTx).Methods("GET")

	return r
}
