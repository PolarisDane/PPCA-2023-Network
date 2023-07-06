package main

import (
	"fmt"
	"net"
)

func main() {
	LocalIP := net.ParseIP("127.0.0.1")
	lAddr := &net.UDPAddr{IP: LocalIP, Port: 10000}
	rAddr := &net.UDPAddr{IP: LocalIP, Port: 8080}
	tarconn, err := net.DialUDP("udp", lAddr, rAddr)
	//UDP客户端向代理服务器发起代理请求
	fmt.Println("Dialed")
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	var buf[512] byte
	buf[0] = 0x05
	buf[1] = 0x01
	buf[2] = 0x00
	tarconn.Write(buf[:3])
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
	tarconn.Write(buf[:10])
	defer tarconn.Close()
}