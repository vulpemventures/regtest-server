package main

import (
	"log"
	"net/http"

	"github.com/vulpemventures/regtest-server/regtest/router"
)

func main() {
	config, err := generateConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	bitcoinClient := &router.RegTest{}
	err = bitcoinClient.New()
	if err != nil {
		log.Fatal(err)
	}
	defer bitcoinClient.Shutdown()

	liquidClient := &router.Liquid{}
	err = liquidClient.New()
	if err != nil {
		log.Fatal(err)
	}
	/**
	** Start new JSON HTTP/1 router
	 */
	r := router.New(bitcoinClient, liquidClient)

	log.Println("Starting server at " + config.Address + ":" + config.Port)
	if err = http.ListenAndServe(config.Address+":"+config.Port, r); err != nil {
		log.Fatal(err)
	}
}
