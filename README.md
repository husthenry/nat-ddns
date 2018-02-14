# nat-ddns

* 2018-02-15  
增加心跳机制

* 2018-02-09

### 用法
> 编译服务端：  
go build Server.go 
服务端配置文件示例：
```json
//config.json
{
  "port": 9257,//proxy代理端口
  "client_keys": [
    {
      "client_key": "9e38630ca96540e5b8611e2d0347df9f",//客户端1秘钥
      "mapper_port": 8080 //客户端1监听端口
    },
    {
      "client_key": "9e38630ca96540e5b8611e2d0347df9f2",//客户端2秘钥
      "mapper_port": 18080 //客户端2监听端口
    }
  ]
}
``` 

> 服务端启动：   
./Server --config=./config.json

> 编译客户端：  
go build Client.go  
客户端配置文件示例：
```json
//client_config.json
{
  "uid": "wenj91",//用户uid,多客户端支持标识
  "client_key": "9e38630ca96540e5b8611e2d0347df9f",//客户端访问KEY
  "server": "127.0.0.1:9257",//代理服务器ip:port
  "real_server": "127.0.0.1:9090"//实际服务器ip:port
}
```
> 客户端启动：  
./Client --client_config=./client_config.json

* 2018-02-06  
> 完成了代理转发的核心逻辑，基本功能可以做内网穿透  
