package main

import (
	"service/client"
	"fmt"
	"flag"
)

var cs = client.ClientService{}

func main() {
	fmt.Println("client conn to server>>>>>>>>>>>>>>>>>>>>>")
	var clientConfig string
	flag.StringVar(&clientConfig, "client_config", "./client_config.json", "--client_config=./client_config.json")
	flag.Parse()


	cs.ClientInit(clientConfig)
	cs.ClientStart()

	fmt.Println("client conn close<<<<<<<<<<<<<<<<<<<<<<<<<")
}


