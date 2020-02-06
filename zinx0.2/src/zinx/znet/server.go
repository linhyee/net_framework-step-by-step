package znet

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"log"
	"net"
	"net_framework-step-by-step/zinx0.2/src/zinx/ziface"
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
}

// CallBackToClient 定义当前客户端链接的所绑定handle api(hard code,日后优化)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	log.Println("[Conn Handle] CallBackClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf error ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// Start 启动服务器
func (s *Server) Start() {
	go func() {
		log.Printf("[Start] Server Listenner at IP :%s, Port:%d, starting\n", s.IP, s.Port)
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
			//已经与客户端建立链接,处理业务,这里做一个基本的回显服务
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			//启动当前 的链接业务处理
			dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务资源、状态或者已经开辟的链接信息,进行停止或者回收
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 处理一些启动服务器之后的额外业务

	//阻塞主协程
	select {}
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8000,
	}
	return s
}
