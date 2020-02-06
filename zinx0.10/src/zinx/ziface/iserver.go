package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//添加路由
	AddRouter(uint32, IRouter)
	//获取当前连接管理器
	GetConnMgr() IConnManager
	//注册OnConnStart钩子方法
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop钩子方法
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStart方法
	CallOnConnStart(conn IConnection)
	//调用OnConnStop方法
	CallOnConnStop(conn IConnection)
}
