package server

import (
	"bufio"
	"constants"
	"encoding/json"
	"entity"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"log"
	"myproto"
	"net"
	"strconv"
	"util/proxy"
	"util"
)

type ServerService struct {
	sc entity.ServerConfig
}

var count = 0

var scs = GetScsInstance()
var scks = GetScksInstance()

var uss = UserServerService{}

func (ss *ServerService) ServerInit(config string) {
	byts, err := util.ReadFile(config)
	if nil != err {
		log.Println("read server config err:", err)
		panic(err)
		return
	}

	sc := entity.ServerConfig{}
	err = json.Unmarshal(byts, &sc)
	if nil != err {
		log.Println("server config unmarshal err:", err)
		panic(err)
		return
	}
	ss.sc = sc


	//客户端密钥管理
	scks.ServerClientKeyServiceInit(sc)

	//用户请求处理服务初始化
	uss.UserServerInit(sc)
}

func (ss *ServerService) ServerStart() {
	go uss.UserServerStart()

	log.Println("server start>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(ss.sc.Port))
	if nil != err {
		log.Println("listen to port:", ss.sc.Port, "err:", err)
		panic(err)
	}

	log.Println("server start at port:", ss.sc.Port)

	for {
		conn, err := listen.Accept()
		if nil != err {
			log.Println("accept conn err:", err)
			continue
		}

		go ss.serverHandle(conn)
	}
}

func (ss *ServerService) serverHandle(conn net.Conn) {
	dataChan := make(chan myproto.Msg)
	errChan := make(chan error)

	go proxy.ReadWrapper(dataChan, errChan, conn)
	go ss.serverDataProcess(dataChan, errChan, conn)
}

func (ss *ServerService) serverDataProcess(dataChan chan myproto.Msg, errChan chan error, conn net.Conn) {
	for {
		uid, _ := uuid.NewV4()
		select {
		case msg := <-dataChan:
			msgBytes, _ := json.Marshal(msg)
			log.Println("recv data from client:", string(msgBytes))

			clientKey := *msg.Key
			switch int(*msg.MsgType) {
			case constants.MSG_TYPE_AUTH:

				if !scks.IsContainsKey(clientKey) {
					log.Println("this client key is not exists!!! clientKey:", clientKey)
					conn.Close()
					return
				}

				isC := scs.IsContainsChannel(clientKey)
				if isC {
					scs.RemoveChannel(clientKey)
				}
				count++
				channel := entity.Channel{
					Id:   count,
					Key:  clientKey,
					Uri:  uid.String(),
					Conn: conn,
					Writable:true,
					SubChan: make(map[string]entity.Channel),
				}
				scs.AddChannel(channel)

				authMsg := myproto.Msg{
					Id:      proto.Int(count),
					MsgType: proto.Int32(constants.MSG_TYPE_AUTH),
					Key:     proto.String(clientKey),
					Uri:     proto.String(uid.String()),
					Data:    []byte("server_auth"),
				}
				_, err := proxy.MsgWrite(authMsg, conn)
				if nil != err {
					log.Println("server send auth pkg failed!", err)
				}
			case constants.MSG_TYPE_HEATBEAT:
				heatBeatMsg := myproto.Msg{
					Id:      proto.Int(count),
					MsgType: proto.Int32(constants.MSG_TYPE_AUTH),
					Key:     proto.String(clientKey),
					Uri:     proto.String(uid.String()),
					Data:    []byte("pong"),
				}
				_, err := proxy.MsgWrite(heatBeatMsg, conn)
				if nil != err {
					log.Println("server send auth pkg failed!", err)
				}
			case constants.MSG_TYPE_CONNECT:
				//set sub_channel writable to true
				key := *msg.Key
				uri := *msg.Uri
				channel := scs.GetSubChannel(key, uri)
				channel.Writable = true
				scs.AddSubChannel(channel)
				log.Println("key:", key, "uri:", uri, " set writable to true success!")
			case constants.MSG_TYPE_DISCONNECT:
				key := *msg.Key
				uri := *msg.Uri
				channel := scs.GetSubChannel(key, uri)
				channel.Conn.Close()
				scs.RemoveSubChannel(key, uri)
				log.Println("key:", key, "uri:", uri, " disconn success!")
			case constants.MSG_TYPE_TRANS:
				//write to user channel
				key := *msg.Key
				uri := *msg.Uri
				channel := scs.GetSubChannel(key, uri)
				w := bufio.NewWriter(channel.Conn)
				w.Write(msg.Data)
				w.Flush()
			}
		case err := <-errChan:
			if nil != err {
				log.Println("An error occured:", err.Error())
				return
			}
		}
	}
}
