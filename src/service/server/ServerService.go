package server

import (
	"encoding/json"
	"entity"
	"log"
	"net"
	"strconv"
	"util"
)

//代理服务管理
type ServerService struct {
	sc entity.ServerConfig
}

//get server channel service
var scs = GetScsInstance()

//get server client key service
var scks = GetScksInstance()

//get user server service instance
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

	//server channel process

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

		pss := ProxyServerService{true, 0}
		go pss.ServerHandle(conn)
	}
}

