package znet

import "net_framework-step-by-step/zinx0.4/src/zinx/ziface"

//实现router时,先嵌入这个BaseRouter基类,然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

//方法为空,不是所有业务都需处理PreHandle PostHandle

// PreHandle 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 在处理conn业务的主方法
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
