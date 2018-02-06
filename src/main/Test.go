package main

import (
	"errors"
	"fmt"
	"net"
	"bufio"
)

var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")

func main()  {

	b1 := []int{1, 2, 3}
	b2 := []int{4, 5}
	b1 = append(b1, b2...)
	fmt.Println(b1)

	l, _ := net.Listen("tcp", ":9090")
	for{
		conn2, _ := l.Accept()
		conn, err := net.Dial("tcp", "www.baidu.com:80")
		if err != nil {
			// handle error
		}
		//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		//status, err := bufio.NewReader(conn).ReadString('\n')
		//fmt.Println(status)

		//conn, _ := net.Dial("tcp", "www.baidu.com:80")
		//
		w := bufio.NewWriter(conn)
		w.Write([]byte("GET / HTTP/1.1\r\nHost: www.baidu.com\r\n\r\n"))
		w.Flush()
		//conn.Write([]byte("Host: 127.0.0.1:8888 \r\n\r\n"))

		//util.ConnTransfer(conn2, conn)

		var buf []byte
		written := int64(0)
		if buf == nil {
			buf = make([]byte, 32*1024)
		}
		for {
			nr, er := conn.Read(buf)
			if nr > 0 {
				nw, ew := conn2.Write(buf[0:nr])
				fmt.Print(string(buf[:nr]))
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != EOF {
					err = er
				}
				break
			}
		}
		fmt.Println("written byte len:", written)

		//r := bufio.NewReader(conn)
		//buf := make([]byte, 1024)
		//for{
		//	i, _ := r.Read(buf)
		//	if i == 0 {
		//		break
		//	}
		//	fmt.Print(string(buf[:i]))
		//	conn2.Write(buf[0:i])
		//}
	}

	//
	//buf := make([]byte, 14633)
	//fmt.Println(len(buf))

}
