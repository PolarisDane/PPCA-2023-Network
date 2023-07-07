package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
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
	//协商认证
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
	//代理请求
	conn.Read(buf[:])
	taraddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	fromaddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10000")
	tarconn, err := net.DialUDP("udp", fromaddr, taraddr)
	//分配的代理服务器端口为9000
	buf[0] = 0x00
	buf[1] = 0x00
	buf[2] = 0x00
	buf[3] = 0x01
	buf[4] = 0x7f
	buf[5] = 0x00
	buf[6] = 0x00
	buf[7] = 0x01//本机地址
	buf[8] = 0x04
	buf[9] = 0xd2//端口1234
	buf[10] = 0x41//A
	tarconn.Write(buf[:11])
	count, _, err := tarconn.ReadFromUDP(buf[:])
	fmt.Println(string(buf[:count]))
	buf[0] = 0x00
	buf[1] = 0x00
	buf[2] = 0x00
	buf[3] = 0x01
	buf[4] = 0x7f
	buf[5] = 0x00
	buf[6] = 0x00
	buf[7] = 0x01//本机地址
	buf[8] = 0x04
	buf[9] = 0xd2//端口1234
	buf[10] = 0x42//B
	tarconn.Write(buf[:11])
	count, _, err = tarconn.ReadFromUDP(buf[:])
	fmt.Println(string(buf[:count]))
	buf[0] = 0x00
	buf[1] = 0x00
	buf[2] = 0x00
	buf[3] = 0x01
	buf[4] = 0x7f
	buf[5] = 0x00
	buf[6] = 0x00
	buf[7] = 0x01//本机地址
	buf[8] = 0x04
	buf[9] = 0xd2//端口1234
	buf[10] = 0x43//C
	tarconn.Write(buf[:11])
	count, _, err = tarconn.ReadFromUDP(buf[:])
	fmt.Println(string(buf[:count]))
	buf[0] = 0x04
	tarconn.Write(buf[:1])
	defer conn.Close()
	defer tarconn.Close()
}