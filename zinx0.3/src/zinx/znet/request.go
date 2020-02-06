package znet

import "net_framework-step-by-step/zinx0.3/src/zinx/ziface"

type Request struct {
	//已经和客户端建立的链接
	conn ziface.IConnection
	//客户端请求的数据
	data []byte
}

// GetConnection 获取当前链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求的消息数据
func (r *Request) GetData() []byte {
	return r.data
}
