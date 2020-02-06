package main

import (
	"fmt"
	"io"
	"net"
	"net_framework-step-by-step/zinx0.8/src/zinx/znet"
	"time"
)

//模拟客户端
func main() {
	fmt.Println("client1 start....")
	time.Sleep(1 * time.Second)
	//1 直接连接远程服务器,得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("client start error, exit!")
		return
	}
	//2 链接调用write写数据
	for {
		//发送封包的Message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("zinx0.8 client test message 1")))
		if err != nil {
			fmt.Println("pack error ", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		//服务器回包
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read error ", err)
			return
		}
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("unpack error ", err)
			return
		}
		if msgHead.GetMsgLen() > 0 {
			data := make([]byte, msgHead.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("read msg error", err)
				return
			}
			msgHead.SetData(data)

			fmt.Println("recv server msg:id=", msgHead.GetMsgId(), ",data=", string(msgHead.GetData()))
		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
