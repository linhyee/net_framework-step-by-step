package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net_framework-step-by-step/zinx0.7/src/zinx/utils"
	"net_framework-step-by-step/zinx0.7/src/zinx/ziface"
)

//封包,拆包的具体模块
type DataPack struct{}

// NewDataPack 拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节缓冲区
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// UnPack 拆包方法 (将包的Head信息读出来)之后再根据head信息的data的长度,再进
func (db *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个存放bytes字节的缓冲区
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息,得到dataLen和MsgID
	msg := &Message{}
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读取MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too larger msg data recv!")
	}
	return msg, nil
}
