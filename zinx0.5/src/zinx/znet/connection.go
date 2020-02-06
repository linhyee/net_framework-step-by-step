package znet

import (
	"errors"
	"io"
	"log"
	"net"
	"net_framework-step-by-step/zinx0.5/src/zinx/ziface"
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
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg的Head头 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			log.Println("read msg header error ", err)
			break
		}
		//拆包,得到MsgId和MsgDataLen放在Msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			log.Println("unpack msg error ", err)
			break
		}
		//根据dataLen再次读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
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

// SendMsg 发送数据,将据发送远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDataLen|MsgId|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		log.Println("Pack msg error msgId:", msgId)
		return errors.New("Pack msg error")
	}
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		log.Println("")
		return errors.New("conn write error")
	}
	return nil
}
