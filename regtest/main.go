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
	client := &router.RegTest{}
	err = client.New()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Mine some blocks to enable faucet service
	_, err = client.Mine(200)
	if err != nil {
		log.Println("Warning: an error occured when mining blocks, please check the following error and manually mine at least 100 blocks to enable faucet service.\n", err)
	}

	r := router.New(client)

	log.Println("Starting server at " + config.Address + ":" + config.Port)
	if err = http.ListenAndServe(config.Address+":"+config.Port, r); err != nil {
		log.Fatal(err)
	}
}
