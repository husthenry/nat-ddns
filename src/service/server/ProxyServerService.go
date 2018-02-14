package server

import (
	"bufio"
	"constants"
	"encoding/json"
	"entity"
	"github.com/golang/protobuf/proto"
	"log"
	"myproto"
	"net"
	"time"
	"util/proxy"
	"util/uuid"
)

type ProxyServerService struct {
	IsConnected bool
	count       int
}

func (pss *ProxyServerService) serverHandle(conn net.Conn) {
	dataChan := make(chan myproto.Msg)
	errChan := make(chan error)
	heartBeatChan := make(chan myproto.Msg)

	go proxy.ReadWrapper(dataChan, errChan, conn)
	go pss.serverDataProcess(dataChan, heartBeatChan, errChan, conn)
}

func (pss *ProxyServerService) serverDataProcess(dataChan chan myproto.Msg, heartBeatChan chan myproto.Msg,
	errChan chan error, conn net.Conn) {

	for {
		uid := uuid.GetRandomUUID()
		select {
		case msg := <-dataChan:
			msgBytes, _ := json.Marshal(msg)
			log.Println("recv data from client:", string(msgBytes))

			clientKey := *msg.Key
			switch int(*msg.MsgType) {
			case constants.MSG_TYPE_AUTH:

				//判断客户端是否合法
				if !scks.IsContainsKey(clientKey) {
					log.Println("this client key is not exists!!! clientKey:", clientKey)
					conn.Close()
					return
				}

				//如果之前有连接,直接移除之前客户端重新接入
				isC := scs.IsContainsChannel(clientKey)
				if isC {
					tmpChannel := scs.GetChannel(clientKey)
					tmpChannel.Conn.Close()
					scs.RemoveChannel(clientKey)
				}

				//重新接入客户端
				pss.count++
				channel := entity.Channel{
					Id:       pss.count,
					Key:      clientKey,
					Uri:      uid,
					Conn:     conn,
					Writable: true,
					SubChan:  make(map[string]entity.Channel),
				}
				scs.AddChannel(channel)

				authMsg := myproto.Msg{
					Id:      proto.Int(pss.count),
					MsgType: proto.Int32(constants.MSG_TYPE_AUTH),
					Key:     proto.String(clientKey),
					Uri:     proto.String(uid),
					Data:    []byte("server_auth"),
				}
				_, err := proxy.MsgWrite(authMsg, conn)
				if nil != err {
					log.Println("server send auth pkg failed!", err)
				}

				//开启心跳监控
				go pss.ServerHeartBeatProcess(heartBeatChan, channel, 60)

			case constants.MSG_TYPE_HEATBEAT:
				heartBeatChan <- msg
				heatBeatMsg := myproto.Msg{
					Id:      proto.Int(pss.count),
					MsgType: proto.Int32(constants.MSG_TYPE_HEATBEAT),
					Key:     proto.String(clientKey),
					Uri:     proto.String(uid),
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

func (pss *ProxyServerService) ServerHeartBeatProcess(heartBeatChan chan myproto.Msg, channel entity.Channel, timeout int) {
	log.Println("heart beat process start>>>>>>>>>>>>>>>>>>>>>>>>>>>>key:", channel.Key)

	for pss.IsConnected {
		select {
		case heartBeatMsg := <-heartBeatChan:
			log.Println("Key:", *heartBeatMsg.Key, "心跳:", string(heartBeatMsg.Data))
			channel.Conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		case <-time.After(time.Duration(timeout) * time.Second):
			//心跳异常结束客户端链接
			log.Println("Key:", channel.Key, "conn dead now")
			channel.Conn.Close()
			pss.IsConnected = false
			break
		}
	}
	log.Println("heart beat process end<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<key:", channel.Key)
}
