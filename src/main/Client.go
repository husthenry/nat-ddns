package main

import (
	"service/client"
	"flag"
	"log"
)

var cs = client.ClientService{0, false}

func main() {
	log.Println("client conn to server>>>>>>>>>>>>>>>>>>>>>")
	var clientConfig string
	flag.StringVar(&clientConfig, "client_config", "./client_config.json", "--client_config=./client_config.json")
	flag.Parse()


	cs.ClientInit(clientConfig)
	cs.ClientStart()

	log.Println("client conn close<<<<<<<<<<<<<<<<<<<<<<<<<")
}


