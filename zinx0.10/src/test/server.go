package main

import (
	"../zinx/znet"
	"log"
	"net_framework-step-by-step/zinx0.10/src/zinx/ziface"
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

// HandleFunc
func HandleFunc(request ziface.IRequest) {
	log.Println("Call HandleFunc Router...")
	log.Println("recv from client: msgID=", request.GetMsgId(),
		",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(203, []byte("Hello, handle func"))
	if err != nil {
		log.Println(err)
	}
}

func OnConnStart(conn ziface.IConnection) {
	log.Println("Begin connection")
	if err := conn.SendMsg(202, []byte("Wecome!"+conn.RemoteAddr().String())); err != nil {
		log.Println(err)
	}
	//给当前的连接设置一些属性
	log.Println("Set Conn Name, Home...")
	conn.SetProperty("Name", "John")
	conn.SetProperty("Home", "myHome")
}

func OnConnStop(conn ziface.IConnection) {
	log.Println("End connection")
	log.Println("connID=", conn.GetConnID(), " disconnected...")

	//获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		log.Println("Name=", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		log.Println("Home=", home)
	}
}

// 基于zinx框架开发的服务器端应用程序
func main() {
	//1 创建一个server实例,使用Zinx的API
	s := znet.NewServer()
	//2 注册连接Hook钩子方法
	s.SetOnConnStart(OnConnStart)
	s.SetOnConnStop(OnConnStop)
	//3 注册路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.AddRouter(3, znet.HandleFunc(HandleFunc))
	//4 启动server
	s.Serve()
}
