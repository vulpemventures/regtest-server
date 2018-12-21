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

	r.HandleFunc("/send/{address}", r.RegTestClient.SendTo)
	r.HandleFunc("/broadcast/{tx}", r.RegTestClient.Broadcast)
	r.HandleFunc("/utxos/{address}", r.RegTestClient.GetUtxos)
	r.HandleFunc("/fees", r.RegTestClient.EstimateFees)
	r.HandleFunc("/ping", r.RegTestClient.Ping)
	r.HandleFunc("/txs/{hash}", r.RegTestClient.GetTx)

	return r
}
