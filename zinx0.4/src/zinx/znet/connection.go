package znet

import (
	"log"
	"net"
	"net_framework-step-by-step/zinx0.4/src/zinx/utils"
	"net_framework-step-by-step/zinx0.4/src/zinx/ziface"
)

//链接模块
type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	//告知当前链接已经退出的/停止的channel
	ExitChan chan bool
	//该链接处理的Router
	Router ziface.IRouter
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// StartReader 链接的读业务访求
func (c *Connection) StartReader() {
	log.Println("Reader Goroutine is running...")
	defer log.Println("connID=", c.ConnID)
	defer c.Stop()

	for {
		//读取客户端的数据到缓冲区中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			log.Println("recv buf error", err)
			continue
		}
		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//从路由中,找到注册绑定的Conn对应的router调用
		//执行注册路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// Start 启动链接 让当前的链接准备工作
func (c *Connection) Start() {
	log.Println("Conn Start().. ConnID=", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//TODO 启动从当前链接写数据业务
}

// Stop 停止链接 结束当前链接工作
func (c *Connection) Stop() {
	log.Println("Conn Stop().. ConnID=", c.ConnID)
	//如果链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket链接
	_ = c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前链接的绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取远程客户羰的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态 IP端口
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据,将据发送远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
