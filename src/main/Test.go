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

	c1 := entity.ClientKeyConfig{
		ClientKey: "9e38630ca96540e5b8611e2d0347df9f",
		MapperPort: 8080,
	}
	c2 := entity.ClientKeyConfig{
		ClientKey: "9e38630ca96540e5b8611e2d0347df9f2",
		MapperPort: 8081,
	}
	var clientKeyConfig []entity.ClientKeyConfig
	clientKeyConfig = append(clientKeyConfig, c1, c2)
	sc := entity.ServerConfig{
		Port: 9898,
		ClientKeys: clientKeyConfig,
	}
	byts,_ = json.Marshal(sc)

	fmt.Println(string(byts))
}
