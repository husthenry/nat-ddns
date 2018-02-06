package main

import (
	"service/server"
	"fmt"
)

var ss = server.ServerService{}

const port  = 9898
const userPort = 9191

func main() {
	fmt.Println("server start>>>>>>>>>>>>>>>>>>>>>")

	ss.ServerInit(port, userPort)

	ss.ServerStart()

	fmt.Println("server end<<<<<<<<<<<<<<<<<<<<<<<")
}
