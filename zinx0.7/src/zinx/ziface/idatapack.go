package ziface

//封包、拆包、模块
//直接面向TCP连接的数据流,用于处理在TCP之上应用层`粘包`问题
type IDataPack interface {
	//获取包的头长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(msg IMessage) ([]byte, error)
	//拆包方法
	UnPack([]byte) (IMessage, error)
}
