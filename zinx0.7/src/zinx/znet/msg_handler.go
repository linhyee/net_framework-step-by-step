package znet

import (
	"log"
	"net_framework-step-by-step/zinx0.7/src/zinx/ziface"
	"strconv"
)

//消息处理模块的实现
type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度或执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		log.Println("api msgID= ", request.GetMsgId(), " 's handler not found")
		return
	}
	//2 根据MsgID调度对应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		log.Panic("repeat api, msgID=" + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	log.Println("Add api MsgID=", msgID, " successfully")
}
