package entity

type ClientKeyConfig struct {
	ClientKey  string `json:"client_key"`
	MapperPort int    `json:"mapper_port"`
}

type ServerConfig struct {
	Port       int                 `json:"port"`
	ClientKeys [](ClientKeyConfig) `json:"client_keys"`
}
