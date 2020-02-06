package ziface

// IRequest接口
// 实际上是把客户端请求信息的链接信息,和请求的数据包封装为Request
type IRequest interface {
	//获取当前链接
	GetConnection() IConnection
	//获取请求的消息数据
	GetData() []byte
}
