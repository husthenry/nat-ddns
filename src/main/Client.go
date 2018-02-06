package main

import (
	"service/client"
	"fmt"
)

var cs = client.ClientService{}

const(
	CLIENT_KEY = "MY_KEY"
	SERVER_ADDR = ":9898"
)

func main() {
	fmt.Println("client conn to server>>>>>>>>>>>>>>>>>>>>>")

	cs.ClientInit(SERVER_ADDR, CLIENT_KEY)
	cs.ClientStart()

	fmt.Println("client conn close<<<<<<<<<<<<<<<<<<<<<<<<<")
}


