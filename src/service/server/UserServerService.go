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
//这里对于tcp的处理方案只能同过端口来进行处理 已经添加done
//对于http可以同过url或者上下文根路径进行处理
var clientKeyConfig []entity.ClientKeyConfig

func (sucs *UserServerService) UserServerInit(sc entity.ServerConfig) {
	sucs.sc = sc
	for i:=0; i<len(sc.ClientKeys); i++  {
		clientKeyConfig = append(clientKeyConfig, sc.ClientKeys[i])
	}
}

func (sucs *UserServerService) UserServerStart() {
	log.Println("user server start>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	for i:=0; i<len(clientKeyConfig); i++ {
		go sucs.userServerStart(clientKeyConfig[i])
	}
}

func (sucs *UserServerService) userServerStart(ckc entity.ClientKeyConfig) {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(ckc.MapperPort))
	if nil != err {
		log.Println("listen to port:", ckc.MapperPort, "err:", err)
		panic(err)
	}

	log.Println("user server start at port:", ckc.MapperPort, "client key:", ckc.ClientKey)

	for {
		conn, err := listen.Accept()
		if nil != err {
			log.Println("accept conn err:", err)
			continue
		}

		go sucs.userServerHandle(conn, ckc)
	}
}

func (sucs *UserServerService) userServerHandle(conn net.Conn, ckc entity.ClientKeyConfig) {
	isExistsProxyChan := uscs.IsContainsChannel(ckc.ClientKey)
	if !isExistsProxyChan {
		log.Println("no proxy chann exists:", ckc.ClientKey)
		conn.Close()
		return
	}

	dataChan := make(chan myproto.Msg)
	errChan := make(chan error)
	connCount++

	//user_channel request ur
	channel := entity.Channel{
		Id:       connCount,
		Key:      ckc.ClientKey,
		Uri:      uuid.GetRandomUUID(),
		Writable: false,
		Conn:     conn,
		SubChan:  make(map[string]entity.Channel),
	}

	//添加子通道
	uscs.AddSubChannel(channel)

	//发送连接通知给客户端，等待响应
	isConn := userHandleConn(channel, ckc)
	if !isConn {
		log.Println("conn to client failed!")
		return
	}

	//连接响应后处理用户请求
	go sucs.userRequestWrapper(dataChan, errChan, conn, channel, ckc)
	go sucs.userDataProcess(dataChan, errChan, conn, ckc)
}

func userHandleConn(channel entity.Channel, ckc entity.ClientKeyConfig) bool {
	//发送连接通知给客户端，等待响应
	connMsg := myproto.Msg{
		Id:      proto.Int(connCount),
		MsgType: proto.Int(constants.MSG_TYPE_CONNECT),
		Key:     proto.String(ckc.ClientKey),
		Uri:     proto.String(channel.Uri),
		Data:    []byte("user_conn"),
	}

	proxyChan := uscs.GetChannel(ckc.ClientKey)

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

func (sucs *UserServerService) userRequestWrapper(dataChan chan myproto.Msg, errChan chan error,
	conn net.Conn, channel entity.Channel, ckc entity.ClientKeyConfig) {
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
			Key:     proto.String(ckc.ClientKey),
			Uri:     proto.String(channel.Uri),
			Data:    buf[:i],
		}

		dataChan <- msg
	}
}

func (sucs *UserServerService) userDataProcess(dataChan chan myproto.Msg, errChan chan error, conn net.Conn,
	ckc entity.ClientKeyConfig) {
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
				channel := uscs.GetChannel(ckc.ClientKey)
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
