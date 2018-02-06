package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")

func main()  {
	var flag = false
	go func() {
		time.Sleep(15*time.Second)
		flag = true
	}()

	for  {
		if flag && true {
			fmt.Println("timeout")
			break
		}
	}
}
