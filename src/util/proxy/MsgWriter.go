package proxy

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"myproto"
	"net"
	"constants"
	"util/math"
)


func MsgWrite(msg myproto.Msg, conn net.Conn) (int, error){
	dataBuf, err := proto.Marshal(&msg)
	if nil != err {
		fmt.Println("msg marshal err:", err)
		return 0, constants.ErrMarshal
	}

	dataLen := len(dataBuf)

	w := bufio.NewWriter(conn)

	frame := append(math.IntToBytes(dataLen), dataBuf...)
	wc, err := w.Write(frame)
	w.Flush()
	if nil != err {
		return wc, constants.ErrShortWrite
	}

	//bts,_:=json.Marshal(msg)
	//log.Println("msg:", string(msg.Data))
	//log.Println("msg write success:", string(bts))

	return wc, nil
}
