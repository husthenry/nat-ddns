package server

import (
	"bufio"
	"constants"
	"entity"
	"log"
	"github.com/golang/protobuf/proto"
	"myproto"
	"net"
	"strconv"
	"util/uuid"
	"util/proxy"
)

type UserServerService struct {
}

var up int
var connCount = 0
var uscs = GetScsInstance()

//var uscks = ServerClientKeyService{}

const (
	KEY = "MY_KEY"
)

func (sucs *UserServerService) UserServerInit(port int) {
	up = port
}

func (sucs *UserServerService) UserServerStart() {
	log.Println("user server start>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(up))
	if nil != err {
		log.Println("listen to port:", up, "err:", err)
		panic(err)
	}

	log.Println("user server start at port:", up)

	for {
		conn, err := listen.Accept()
		if nil != err {
			log.Println("accept conn err:", err)
			continue
		}

		go sucs.userServerHandle(conn)
	}
}

func (sucs *UserServerService) userServerHandle(conn net.Conn) {
	dataChan := make(chan myproto.Msg)
	errChan := make(chan error)
	connCount++

	//user_channel request ur
	channel := entity.Channel{
		Id:   connCount,
		Key:  KEY,
		Uri:  uuid.GetRandomUUID(),
		Conn: conn,
	}

	//添加子通道
	uscs.AddSubChannel(channel)

	//发送连接通知给客户端，等待响应
	userHandleConn(channel)

	//连接响应后处理用户请求
	go sucs.userRequestWrapper(dataChan, errChan, conn, channel)
	go sucs.userDataProcess(dataChan, errChan, conn)
}

func userHandleConn(channel entity.Channel){
	//发送连接通知给客户端，等待响应
	connMsg := myproto.Msg{
		Id: proto.Int(connCount),
		MsgType:proto.Int(constants.MSG_TYPE_CONNECT),
		Key:proto.String(KEY),
		Uri: proto.String(channel.Uri),
		Data: []byte("user_conn"),
	}

	proxyChan := uscs.GetChannel(KEY)

	_, err := proxy.MsgWrite(connMsg, proxyChan.Conn)
	if nil != err {
		log.Println("send conn to client err:", err)
		uscs.GetSubChannel(channel.Key, channel.Uri).Conn.Close()
		uscs.RemoveSubChannel(channel.Key, channel.Uri)
		return
	}

}

func (sucs *UserServerService) userRequestWrapper(dataChan chan myproto.Msg, errChan chan error, conn net.Conn, channel entity.Channel) {
	for {
		log.Println("start read user request>>>>>>>>>>>>>")
		r := bufio.NewReader(conn)
		buf := make([]byte, 1024)
		i, err := r.Read(buf)
		if nil != err {
			errChan <- err
			return
		}

		msg := myproto.Msg{
			Id:      proto.Int(connCount),
			MsgType: proto.Int(constants.MSG_TYPE_TRANS),
			Key:     proto.String(KEY),
			Uri:     proto.String(channel.Uri),
			Data:    buf[:i],
		}

		dataChan <- msg
	}
}

func (sucs *UserServerService) userDataProcess(dataChan chan myproto.Msg, errChan chan error, conn net.Conn) {
	for {
		select {
		case msg := <-dataChan:
			log.Println("msg from client:")
			log.Println(" msgId:", *msg.Id, ", msgType:", *msg.MsgType, ", uri:", *msg.Uri,
				", data:", string(msg.Data))
			switch int(*msg.MsgType) {
			case constants.MSG_TYPE_TRANS:
				//write to client
				//直接将请求写入客户端中
				channel := uscs.GetChannel(KEY)
				_, err := proxy.MsgWrite(msg, channel.Conn)
				if nil != err {
					log.Println("msg write err:", err.Error())
				}
			}
		case err := <-errChan:
			if nil != err {
				log.Println("An error occured:", err.Error())
				return
			}
		}
	}
}
