package znet

import (
	"log"
	"net_framework-step-by-step/zinx0.10/src/zinx/utils"
	"net_framework-step-by-step/zinx0.10/src/zinx/ziface"
	"strconv"
)

//消息处理模块的实现
type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的数量
	WorkerPoolSize uint32
}

// NewMsgHandle 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// StartWorkerPool 启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize分别开启Worker，每个Worker用一个go来启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker启动
		//1 当前的worker对应的chan消息队列初始化
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的Worker,阻塞等待消息从channel传进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {
	log.Println("WorkID=", workID, " is started")
	//不断阻塞等待对应消息队列的消息
	for {
		select {
		//处理消息,执行消息当前绑定的任务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 向工作池投递消息
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 使用平均分配算法
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	log.Println("Add connID=", request.GetConnection().GetConnID(),
		" msgID=", request.GetMsgId(), " to Worker[workID=", workerID, "]")

	//2 投递消息
	mh.TaskQueue[workerID] <- request
}
