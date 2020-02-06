package main

import (
	"fmt"
	"net"
	"time"
)

//模拟客户端
func main() {
	fmt.Println("client start....")
	time.Sleep(1 * time.Second)
	//1 直接连接远程服务器,得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("client start error, exit!")
		return
	}
	//2 链接调用write写数据
	for {
		_, err := conn.Write([]byte("Hello, zinx v0.3"))
		if err != nil {
			fmt.Println("write conn error", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error", err)
			return
		}
		fmt.Printf("server call back: %s, cnt=%d\n", buf[:cnt], cnt)

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
