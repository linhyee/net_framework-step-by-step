package znet

import (
	"errors"
	"io"
	"log"
	"net"
	"net_framework-step-by-step/zinx0.9/src/zinx/utils"
	"net_framework-step-by-step/zinx0.9/src/zinx/ziface"
)

//链接模块
type Connection struct {
	//当前Conn属于哪个Server(分布式,扩展多server)
	TcpServer ziface.IServer
	//当前链接的socket TCP套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	//告知当前链接已经退出的/停止的channel(由reader告知writer,让writer退出)
	ExitChan chan bool
	//无缓冲管道,用于读、写goroutine之间的消息通信
	msgChan chan []byte
	//消息的管理MsgID和对应的处理业务API
	MsgHandler ziface.IMsgHandle
}

// NewConnection 初始化链接模块的方法
func NewConnection(srv ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  srv,
		Conn:       conn,
		ConnID:     connID,
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}
	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// StartReader 链接的读业务访求
func (c *Connection) StartReader() {
	log.Println("[Reader] Goroutine is running...")
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制,将消息发送给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中,找到注册绑定的Conn对应的router调用
			//根据绑定好的MsgID找到对应处理api业务
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 写消息Goroutine,用户将消息发送给客户端的模块
func (c *Connection) StartWriter() {
	log.Println("[Writer] goroutine is running")
	defer log.Println(c.RemoteAddr().String(), "conn writer exit")

	//不断的阻塞地等待channnel的消息,进写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据写
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出,此时writer返回
			return
		}
	}
}

// Start 启动链接 让当前的链接准备工作
func (c *Connection) Start() {
	log.Println("Conn Start().. ConnID=", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//启动从当前链接写数据业务
	go c.StartWriter()
	//按照开发者传递进来的hook方法,执行OnStart方法
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止链接 结束当前链接工作
func (c *Connection) Stop() {
	log.Println("Conn Stop().. ConnID=", c.ConnID)
	//如果链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//调用开发都注册的HookFn OnConnStop
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	_ = c.Conn.Close()

	//将当前连接从connMgr中删掉
	c.TcpServer.GetConnMgr().Remove(c)

	//告知writer关闭
	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
		return errors.New("connection closed when send msg")
	}
	//将data进行封包 MsgDataLen|MsgId|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		log.Println("Pack msg error msgId:", msgId)
		return errors.New("pack msg error")
	}
	//将数据发送给msgChan
	c.msgChan <- binaryMsg

	return nil
}
