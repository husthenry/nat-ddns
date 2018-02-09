package main

import (
	"entity"
	"encoding/json"
	"fmt"
)

func main()  {
	cc := entity.ClientConfig{
		Uid        : "9e38630ca96540e5b8611e2d0347df9f",
		ClientKey  : "9e38630ca96540e5b8611e2d0347df9f",
		Server     : "127.0.0.1:9898",
		RealServer : "127.0.0.1:9191",
	}
	byts,_ := json.Marshal(cc)

	fmt.Println(string(byts))

	cMap := make(map[string]string)
	cMap["client_key_9e38630ca96540e5b8611e2d0347df9f"]="9e38630ca96540e5b8611e2d0347df9f"
	cMap["client_key_9e38630ca96540e5b8611e2d0347df9f2"]="9e38630ca96540e5b8611e2d0347df9f2"
	sc := entity.ServerConfig{
		Port: 9898,
		UserPort: 9191,
		ClientKey: []map[string]string{cMap},
	}
	byts,_ = json.Marshal(sc)

	fmt.Println(string(byts))
}
