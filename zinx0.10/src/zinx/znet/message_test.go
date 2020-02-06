package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//测试DataPack拆包封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟服务器
	//1.创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("server listen error ", err)
	}
	//创建一个go承载客户端处理业务
	go func() {
		//2 从客户端读取数据,拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error ", err)
				continue
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				//拆包过程
				db := NewDataPack()
				for {
					//1 第一次从conn读, 把包头读出来
					headData := make([]byte, db.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					msgHead, err := db.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//2 第二次从conn读, 根据包头信息读出包内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//根据dataLen长茺再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error ", err)
							return
						}

						//完整的一个消息已经读取完毕
						fmt.Println("recv msgId: ", msg.Id, ", dataLen: ", msg.DataLen, ", data: ", string(msg.Data))
					}

				}
			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("client dial error ", err)
	}
	//创建一个封包对象 db
	db := NewDataPack()

	//模拟分包
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := db.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error ", err)
		return
	}

	//封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := db.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error ", err)
		return
	}
	//将两个包连在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性发送给服务端
	_, _ = conn.Write(sendData1)

	select {}
}
