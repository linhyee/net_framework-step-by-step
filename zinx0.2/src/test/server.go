package main

import (
	"../zinx/znet"
)

// 基于zinx框架开发的服务器端应用程序
func main() {
	//1 创建一个server实例,使用Zinx的API
	s := znet.NewServer("[zinx v0.2]")
	//2 启动server
	s.Serve()
}
