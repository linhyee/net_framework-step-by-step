package main

import (
	"../zinx/znet"
	"log"
	"net_framework-step-by-step/zinx0.7/src/zinx/ziface"
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

// Hello router
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Hello Handle
func (h *HelloZinxRouter) Handle(request ziface.IRequest) {
	log.Println("Call Hello Router....")
	log.Println("recv form client: msgID=", request.GetMsgId(),
		",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello, welcome!"))
	if err != nil {
		log.Println(err)
	}
}

// 基于zinx框架开发的服务器端应用程序
func main() {
	//1 创建一个server实例,使用Zinx的API
	s := znet.NewServer("[zinx v0.7]")
	//2 注册路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//3 启动server
	s.Serve()
}
