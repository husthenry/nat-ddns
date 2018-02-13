package main

import (
	"net"
	"util/proxy"
)

func main()  {

	l,_:= net.Listen("tcp", ":9999")
	for  {
		conn,_:= l.Accept()
		go handleTest(conn)
	}
}

func handleTest(conn net.Conn)  {
	target,_:=net.Dial("tcp", ":22")
	proxy.ConnTransfer(conn, target)
}