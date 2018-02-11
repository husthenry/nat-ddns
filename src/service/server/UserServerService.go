package server

import (
	"bufio"
	"constants"
	"entity"
	"github.com/golang/protobuf/proto"
	"log"
	"myproto"
	"net"
	"strconv"
	"time"
	"util/proxy"
	"util/uuid"
)

type UserServerService struct {
	sc entity.ServerConfig
}

var connCount = 0
var uscs = GetScsInstance()

//todo: 多客户端用户请求管理
//这里k只是单客户端用户请求管理
//这里对于tcp的处理方案只能同过端口来进行处理
//对于http可以同过url或者上下文根路径进行处理
var k string

func (sucs *UserServerService) UserServerInit(sc entity.ServerConfig) {
	sucs.sc = sc
	for _, v := range sc.ClientKey[0]{
		k = v
	}
}

func (sucs *UserServerService) UserServerStart() {
	log.Println("user server start>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(sucs.sc.UserPort))
	if nil != err {
		log.Println("listen to port:", sucs.sc.UserPort, "err:", err)
		panic(err)
	}

	log.Println("user server start at port:", sucs.sc.UserPort)

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
	isExistsProxyChan := uscs.IsContainsChannel(k)
	if !isExistsProxyChan {
		log.Println("no proxy chann exists:", k)
		conn.Close()
		return
	}

	dataChan := make(chan myproto.Msg)
	errChan := make(chan error)
	connCount++

	//user_channel request ur
	channel := entity.Channel{
		Id:       connCount,
		Key:      k,
		Uri:      uuid.GetRandomUUID(),
		Writable: false,
		Conn:     conn,
		SubChan:  make(map[string]entity.Channel),
	}

	//添加子通道
	uscs.AddSubChannel(channel)

	//发送连接通知给客户端，等待响应
	isConn := userHandleConn(channel)
	if !isConn {
		log.Println("conn to client failed!")
		return
	}

	//连接响应后处理用户请求
	go sucs.userRequestWrapper(dataChan, errChan, conn, channel)
	go sucs.userDataProcess(dataChan, errChan, conn)
}

func userHandleConn(channel entity.Channel) bool {
	//发送连接通知给客户端，等待响应
	connMsg := myproto.Msg{
		Id:      proto.Int(connCount),
		MsgType: proto.Int(constants.MSG_TYPE_CONNECT),
		Key:     proto.String(k),
		Uri:     proto.String(channel.Uri),
		Data:    []byte("user_conn"),
	}

	proxyChan := uscs.GetChannel(k)

	_, err := proxy.MsgWrite(connMsg, proxyChan.Conn)
	if nil != err {
		log.Println("send conn to client err:", err)
		uscs.GetSubChannel(channel.Key, channel.Uri).Conn.Close()
		uscs.RemoveSubChannel(channel.Key, channel.Uri)
		return false
	}

	//等待60s客户端响应,如果正常响应则放行，反之断开连接
	var flag = false
	go func() {
		time.Sleep(60 * time.Second)
		flag = true
		log.Println("timeout task finish")
	}()

	for {
		writable := uscs.GetSubChannel(channel.Key, channel.Uri).Writable
		if writable {
			return true
		}
		if flag && !channel.Writable {
			log.Println("wait conn to client timeout!!!")
			uscs.GetSubChannel(channel.Key, channel.Uri).Conn.Close()
			uscs.RemoveSubChannel(channel.Key, channel.Uri)
			return false
		}
	}

	return false

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
			Key:     proto.String(k),
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
				channel := uscs.GetChannel(k)
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
