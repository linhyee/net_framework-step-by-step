package main

import (
	"../zinx/znet"
	"log"
	"net_framework-step-by-step/zinx0.5/src/zinx/ziface"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	log.Println("Call Router Handle...")
	log.Println("recv from client: msgId=", request.GetMsgId(),
		",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		log.Println(err)
	}
}

// 基于zinx框架开发的服务器端应用程序
func main() {
	//1 创建一个server实例,使用Zinx的API
	s := znet.NewServer("[zinx v0.5]")
	//2 注册路由
	s.AddRouter(&PingRouter{})
	//3 启动server
	s.Serve()
}
