package ziface

import (
	"net"
)

//定义链接模块的抽象层
type IConnection interface {
	//启动链接 让当前的链接准备开始工作
	Start()
	//停止链接 结束当前连接工作
	Stop()
	//获取当前链接的绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前连接的模块的链接ID
	GetConnID() uint32
	//获取远程客户端的TCP状态 IP端口
	RemoteAddr() net.Addr
	//发送数据, 将数据发送给远程的客户端
	SendMsg(uint32, []byte) error
}
