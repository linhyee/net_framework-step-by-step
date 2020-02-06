package main

import (
	"../zinx/znet"
	"log"
	"net_framework-step-by-step/zinx0.4/src/zinx/ziface"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test PreHandle
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	log.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		log.Println("call back before ping error")
	}
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	log.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..."))
	if err != nil {
		log.Println("call back ping...ping...ping...")
	}
}

// Test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	log.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping"))
	if err != nil {
		log.Println("call back after ping error")
	}
}

// 基于zinx框架开发的服务器端应用程序
func main() {
	//1 创建一个server实例,使用Zinx的API
	s := znet.NewServer("[zinx v0.4]")
	//2 注册路由
	s.AddRouter(&PingRouter{})
	//3 启动server
	s.Serve()
}
