package znet

import (
	"fmt"
	"log"
	"net"
	"net_framework-step-by-step/zinx0.9/src/zinx/utils"
	"net_framework-step-by-step/zinx0.9/src/zinx/ziface"
)

//IServer的接口实现,定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息管理模块
	MsgHandler ziface.IMsgHandle
	//server的连接管理器
	ConnMgr ziface.IConnManager
	//Hook钩子方法 -连接开始时调用
	OnConnStart func(conn ziface.IConnection)
	//Hook钩子方法 -连接结束时调用
	OnConnStop func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {
	log.Printf("[Zinx] Server Name : %s, listener at IP : %s, Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	log.Printf("[Zinx] Version %s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	log.Printf("[Start] Server Listenner at IP :%s, Port:%d, starting\n", s.IP, s.Port)
	go func() {
		//0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Println("resolve tcp addr error:", err)
			return
		}
		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Println("listen ", s.IPVersion, " error ", err)
			return
		}
		log.Println("start zinx server successfully,", s.Name)

		var cid uint32
		cid = 0

		//3 阻塞等待客户端链拉,处理客户端链接业务
		for {
			//如果有客户端链接, 阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Println("Accept error ", err)
				continue
			}
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端相应一个超出最大连接的错误提示
				log.Println("too many connections[MaxConn=", utils.GlobalObject.MaxConn, "]")
				_ = conn.Close()
				continue
			}
			//已经与客户端建立链接,处理业务,这里做一个基本的回显服务
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前 的链接业务处理
			dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	//将一些服务资源、状态或者已经开辟的链接信息,进行停止或者回收
	log.Println("[STOP]zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 处理一些启动服务器之后的额外业务

	//阻塞主协程
	select {}
}

// AddRouter 添加路由
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	log.Println("Add Router successfully")
}

// GetConnMgr 获
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 注册OnConnStart钩子方法
func (s *Server) SetOnConnStart(hookFn func(conn ziface.IConnection)) {
	s.OnConnStart = hookFn
}

// SetOnConnStop 注册OnConnStop钩子方法
func (s *Server) SetOnConnStop(hookFn func(conn ziface.IConnection)) {
	s.OnConnStop = hookFn
}

// CallOnConnStart 调用OnConnStart方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		log.Println("Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnStop方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		log.Println("Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
