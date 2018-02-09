package entity

type ClientConfig struct {
	Uid        string `json:"uid"`
	ClientKey  string `json:"client_key"`
	Server     string `json:"server"`
	RealServer string `json:"real_server"`
}
