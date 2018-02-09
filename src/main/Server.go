package main

import (
	"service/server"
	"fmt"
	"flag"
)

var ss = server.ServerService{}


func main() {
	fmt.Println("server start>>>>>>>>>>>>>>>>>>>>>")

	var config string
	flag.StringVar(&config, "config", "./config.json", "--config ./config.json")

	flag.Parse()

	ss.ServerInit(config)

	ss.ServerStart()

	fmt.Println("server end<<<<<<<<<<<<<<<<<<<<<<<")
}
