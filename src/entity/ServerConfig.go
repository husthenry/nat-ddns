package entity

type ServerConfig struct {
	Port      int                   `json:"port"`
	UserPort  int                   `json:"user_port"`
	ClientKey [](map[string]string) `json:"client_keys"`
}
