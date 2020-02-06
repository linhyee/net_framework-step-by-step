package znet

import "net_framework-step-by-step/zinx0.10/src/zinx/ziface"

type Request struct {
	//已经和客户端建立的链接
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}

// GetConnection 获取当前链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 获取请求的消息ID
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
