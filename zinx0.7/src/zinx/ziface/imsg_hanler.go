package ziface

//消息管理抽象层
type IMsgHandle interface {
	//调度/执行对应的路由器消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具本的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
