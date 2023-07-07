package main

import (
	"fmt"
	"net"
)

func main() {
	LocalIP := net.ParseIP("127.0.0.1")
	var addr = "127.0.0.1:8080"
	conn, err := net.Dial("tcp", addr)
	//UDP客户端向代理服务器发起代理请求，TCP连接
	fmt.Println("Dialed")
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	var buf[512] byte
	buf[0] = 0x05
	buf[1] = 0x01
	buf[2] = 0x00
	conn.Write(buf[:3])
	conn.Read(buf[:])
	buf[0] = 0x05
	buf[1] = 0x03
	buf[2] = 0x00
	buf[3] = 0x01
	buf[4] = 0x7f
	buf[5] = 0x00
	buf[6] = 0x00
	buf[7] = 0x01
	buf[8] = 0x27
	buf[9] = 0x10
	conn.Write(buf[:10])
	conn.Read(buf[:])
	buf[0] = 0x00
	buf[1] = 0x00
	buf[2] = 0x00
	buf[3] = 0x01
	buf[4] = 0x7f
	buf[5] = 0x00
	buf[6] = 0x00
	buf[7] = 0x01
	buf[8] = 0x04
	buf[9] = 0xd2
	buf[10] = 0x41//A
	conn.Write(buf[:10])
	defer conn.Close()
}