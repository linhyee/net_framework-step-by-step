package znet

import (
	"errors"
	"log"
	"net_framework-step-by-step/zinx0.9/src/zinx/ziface"
	"strconv"
	"sync"
)

//连接管理模块
type ConnManager struct {
	//管理的连接集合
	connections map[uint32]ziface.IConnection
	//保护连接集合的读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建连接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//Add 添加连接
func (mgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map, 加写锁
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()

	//将conn加入到connManager中
	mgr.connections[conn.GetConnID()] = conn
	log.Println("connID=", conn.GetConnID(), "add to ConnManager successfully")
}

//Remove 删除连接
func (mgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map, 加写锁
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	//删除连接信息
	delete(mgr.connections, conn.GetConnID())
	log.Println("connID=", conn.GetConnID(), " removed successfully")
}

//Get 根据connID获取连接
func (mgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map,加读锁
	mgr.connLock.RLock()
	defer mgr.connLock.RUnlock()

	if conn, ok := mgr.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connID=" + strconv.Itoa(int(connID)) + " not found")
}

//Len 得到当前连接总数
func (mgr *ConnManager) Len() int {
	return len(mgr.connections)
}

//ClearConn 清除并终止所有连接
func (mgr *ConnManager) ClearConn() {
	//保护共享资源map,加写锁
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()

	//删除conn连接并停止conn的工作
	for connID, conn := range mgr.connections {
		//停止连接
		conn.Stop()
		//删除
		delete(mgr.connections, connID)
	}
	log.Println("clear all connections successfully, conn num=", mgr.Len())
}
