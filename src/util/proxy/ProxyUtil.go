package proxy

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"io"
	"myproto"
	"net"
	"sync"
	"util/math"
	"log"
)

//代理转发逻辑
func ConnTransfer(currentConn, targetConn net.Conn) {
	var wg sync.WaitGroup
	// 将当前请求转发至目标连接
	go func(toConn net.Conn, conn net.Conn) {
		defer wg.Done()
		wg.Add(1)
		io.Copy(toConn, conn)
		//conn.Close()
	}(targetConn, currentConn)

	// 将目标响应返回给客户端
	go func(conn net.Conn, toConn net.Conn) {
		defer wg.Done()
		wg.Add(1)
		io.Copy(conn, toConn)
		//toConn.Close()
	}(currentConn, targetConn)

	//等待相应的线程都执行完毕
	wg.Wait()
}

//消息读取解码
func ReadWrapper(dataChan chan myproto.Msg, errChan chan error, conn net.Conn) {

	// refer io.Copy
	for {
		//验证conn是否有效
		r := bufio.NewReader(conn)
		lenBuf := make([]byte, 8)
		_, err := r.Read(lenBuf)
		if nil != err {
			log.Println("read from conn err:", err)
			errChan<-err
			break
		}
		dataLen := math.BytesToInt(lenBuf)
		if dataLen > 0 {
			dataBuf := make([]byte, dataLen)
			readCount := 0
			//循环读取，直到将所有字节读取完毕,这里有可能会有问题(数据有问题的时候会引发死循环)
			for{
				if readCount == dataLen {
					break
				}
				n, err := r.Read(dataBuf[readCount:])
				readCount += n
				if nil != err {
					errChan <- err
					break
				}
			}

			//fmt.Println("dataLen:", dataLen, "recvLen:", readCount)
			if readCount == dataLen {
				msg := myproto.Msg{}
				proto.Unmarshal(dataBuf, &msg)
				dataChan <- msg
			}
		}
	}
}
