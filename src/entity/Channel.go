package entity

import (
	"net"
)

type Channel struct {
	Id       int
	Key      string             //User Channel Key
	Uri      string             //User request uri
	Conn     net.Conn           //User request conn
	Writable bool               //channel is writable
	SubChan  map[string]Channel //sub channel
}
