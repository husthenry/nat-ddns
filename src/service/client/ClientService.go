package client

import (
	"bufio"
	"constants"
	"encoding/json"
	"entity"
	"github.com/golang/protobuf/proto"
	"log"
	"myproto"
	"net"
	"strconv"
	"time"
	"util"
	"util/math"
	"util/proxy"
	"util/uuid"
)

type ClientService struct {
	Id          int
	IsConnected bool
}

var cc = entity.ClientConfig{}

// get client channel service
var ccs = GetCcsInstance()

func (cs *ClientService) ClientInit(clientConfig string) {

	byts, err := util.ReadFile(clientConfig)
	if nil != err {
		log.Println("read client config err:", err)
		panic(err)
		return
	}

	err = json.Unmarshal(byts, &cc)
	if nil != err {
		log.Println("unmarshal client config err:", err)
		panic(err)
		return
	}
}

func (cs *ClientService) ClientStart() {

	for {
		if !cs.IsConnected {
			log.Println("conn to server..............")
			conn, err := net.Dial("tcp", cc.Server)
			if nil != err {
				log.Println("dial to server:", cc.Server, " err:", err, " 5s later will be reconn......")
				time.Sleep(5 * time.Second)
				continue
			}

			cs.IsConnected = true

			go cs.clientHandle(conn)
		}else{
			//阻塞,防止cpu被过度利用
			time.Sleep(10*time.Second)
		}
	}
}

func (cs *ClientService) clientHandle(conn net.Conn) {
	cs.cilentAuth(conn)
}

func (cs *ClientService) cilentAuth(conn net.Conn) {
	authMsg := myproto.Msg{
		Id:      proto.Int(cs.Id),
		MsgType: proto.Int32(constants.MSG_TYPE_AUTH),
		Key:     proto.String(cc.ClientKey),
		Uri:     proto.String(uuid.GetRandomUUID()),
		Data:    []byte("client_auth"),
	}

	_, err := proxy.MsgWrite(authMsg, conn)
	if nil != err {
		log.Println("client send auth err")
		cs.IsConnected = false
		return
	}

	r := bufio.NewReader(conn)
	for {
		lenBuf := make([]byte, 8)
		r.Read(lenBuf)
		dataLen := math.BytesToInt(lenBuf)
		if dataLen > 0 {
			dataBuf := make([]byte, dataLen)
			n, err := r.Read(dataBuf)
			if nil != err {
				log.Println("read buffer from server err:", err)
				panic(err)
				break
			}
			if n == dataLen {
				msg := myproto.Msg{}
				proto.Unmarshal(dataBuf, &msg)
				msgBytes, _ := json.Marshal(msg)
				log.Println("recv data from server:", string(msgBytes))
				if int(*msg.MsgType) == constants.MSG_TYPE_AUTH {
					cs.Id = int(*msg.Id)

					//如果存在之前连接,则清除之前连接
					if ccs.IsContainsChannel(*msg.Key) {
						ccs.GetChannel(*msg.Key).Conn.Close()
						ccs.RemoveChannel(*msg.Key)
					}

					//add proxy chan
					proxyChan := entity.Channel{
						Id:       int(*msg.Id),
						Key:      *msg.Key,
						Uri:      *msg.Uri,
						Writable: true,
						Conn:     conn,
						SubChan:  make(map[string]entity.Channel),
					}
					ccs.AddChannel(proxyChan)

					// ping
					go cs.ping(conn)

					dataChan := make(chan myproto.Msg)
					errChan := make(chan error)
					heartBeatChan := make(chan myproto.Msg)
					go cs.ClientHeartBeatProcess(heartBeatChan, proxyChan, 60)
					go proxy.ReadWrapper(dataChan, errChan, conn)
					go cs.clientDataProcess(dataChan, heartBeatChan, errChan, conn)

					cs.IsConnected = true
					break
				}
			}
		}
	}
}

func (cs *ClientService) ping(conn net.Conn) {
	log.Println("send heart beat packet......")

	t := time.NewTicker(30 * time.Second)

	heatBeatCount := 0
	heatBeatErrCount := 0

	for cs.IsConnected {
		select {
		case i := <-t.C:
			heatBeatCount++
			log.Println("ping count:", strconv.Itoa(heatBeatCount), " client ping:",
				i.Format("2006-01-02 15:04:05"))
			authMsg := myproto.Msg{
				Id:      proto.Int(cs.Id),
				MsgType: proto.Int32(constants.MSG_TYPE_HEATBEAT),
				Key:     proto.String(cc.ClientKey),
				Uri:     proto.String(uuid.GetRandomUUID()),
				Data:    []byte(strconv.Itoa(heatBeatCount)),
			}

			_, err := proxy.MsgWrite(authMsg, conn)
			if nil != err {
				heatBeatErrCount++
				log.Println("count:", heatBeatErrCount, "client send heatbeat msg failed!")
			}
		}
	}
}

func (cs *ClientService) clientDataProcess(dataChan chan myproto.Msg, heartBeatChan chan myproto.Msg,
	errChan chan error, conn net.Conn) {
	for {
		select {
		case msg := <-dataChan:
			msgBytes, _ := json.Marshal(msg)
			log.Println("recv data from server:", string(msgBytes))

			switch int(*msg.MsgType) {
			case constants.MSG_TYPE_HEATBEAT:
				//heat beat exception process
				heartBeatChan <- msg

			case constants.MSG_TYPE_CONNECT:
				//handle conn msg: dial to real server
				//dial to real server
				realConn, _ := net.Dial("tcp", cc.RealServer)
				realChan := entity.Channel{
					Id:       int(*msg.Id),
					Key:      *msg.Key,
					Uri:      *msg.Uri,
					Writable: true,
					Conn:     realConn,
					SubChan:  make(map[string]entity.Channel),
				}
				ccs.AddSubChannel(realChan)

				connMsg := myproto.Msg{
					Id:      msg.Id,
					Key:     msg.Key,
					Uri:     msg.Uri,
					MsgType: msg.MsgType,
					Data:    []byte("client_conn_resp"),
				}
				//then resp to proxyChan
				proxy.MsgWrite(connMsg, conn)

				go clientTransDataProcess(msg, realChan)

			case constants.MSG_TYPE_TRANS:
				//handle the trans data
				realChan := ccs.GetSubChannel(*msg.Key, *msg.Uri)
				if realChan.Id == 0 {
					//todo: err process
				}

				//target, _ := net.Dial("tcp", "www.baidu.com:80")
				//str := strings.Replace(string(msg.Data), "127.0.0.1:9191", "www.baidu.com", -1)
				w := bufio.NewWriter(realChan.Conn)
				//log.Println("bytes:", str)
				w.Write(msg.Data)
				w.Flush()
			}
		case err := <-errChan:
			if nil != err {
				log.Println("An error occured:", err.Error())
				cs.IsConnected = false
				return
			}
		}
	}
}

func clientTransDataProcess(msg myproto.Msg, realChan entity.Channel) {
	buf := make([]byte, 32*1024)
	written := int64(0)
	proxyChan := ccs.GetChannel(*msg.Key)
	for {
		i, err := realChan.Conn.Read(buf)
		if i > 0 {
			msg := myproto.Msg{
				Id:      msg.Id,
				MsgType: proto.Int(constants.MSG_TYPE_TRANS),
				Key:     msg.Key,
				Uri:     msg.Uri,
				Data:    buf[:i],
			}
			wc, err2 := proxy.MsgWrite(msg, proxyChan.Conn)
			written += int64(wc)
			if nil != err2 {
				log.Println("Write Error", err2)
			}
		}
		if err != nil {
			if err.Error() != constants.EOF.Error() {
				log.Println("read err:", err)
				//disconn for sub channel
				realChan.Conn.Close()
				ccs.RemoveSubChannel(*msg.Key, *msg.Uri)

				disconnMsg := myproto.Msg{
					Id:      msg.Id,
					MsgType: proto.Int(constants.MSG_TYPE_DISCONNECT),
					Key:     msg.Key,
					Uri:     msg.Uri,
					Data:    []byte(err.Error()),
				}

				proxy.MsgWrite(disconnMsg, proxyChan.Conn)
			}
			break
		}
	}
}

func (cs *ClientService) ClientHeartBeatProcess(heartBeatChan chan myproto.Msg, channel entity.Channel, timeout int) {
	log.Println("heart beat process start>>>>>>>>>>>>>>>>>>>>>>>>>>>>key:", channel.Key)
	for cs.IsConnected {
		select {
		case heartBeatMsg := <-heartBeatChan:
			log.Println("Key:", *heartBeatMsg.Key, "心跳:", string(heartBeatMsg.Data))
			channel.Conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		case <-time.After(time.Duration(timeout) * time.Second):
			//心跳异常结束客户端链接
			log.Println("Key:", channel.Key, "conn dead now")
			channel.Conn.Close()
			cs.IsConnected = false
			break
		}
	}
	log.Println("heart beat process end<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<key:", channel.Key)
}
