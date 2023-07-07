package main

import (
	"fmt"
	"net"
	"encoding/binary"
)

func DistributePort() *net.UDPAddr {
	LocalIP := net.ParseIP("127.0.0.1")
	return &net.UDPAddr{IP: LocalIP, Port: 10000}
	//分配代理服务器端口，未完成
}

func AnswerUDPRequest(conn net.Conn, DistributedAddr *net.UDPAddr) {
	var buf[512] byte
	//未完成
	buf[0] = byte(0x05)
	buf[1] = byte(0x00)
	buf[2] = byte(0x00)
	buf[3] = byte(0x01)
	buf[4] = byte(0x7F)
	buf[5] = byte(0x00)
	buf[6] = byte(0x00)
	buf[7] = byte(0x01)
	buf[8] = byte(0x27)
	buf[9] = byte(0x10)
	conn.Write(buf[:10])
}

func HandleUDP(addr string, conn net.Conn) {
	var buf[512] byte
	//回复代理请求
	clientaddr, err := net.ResolveUDPAddr("udp", addr)
	var DistributedAddr = DistributePort()
	clientconn, err := net.ListenUDP("udp", DistributedAddr)
	AnswerUDPRequest(conn, DistributedAddr)
	count, useraddr, err := clientconn.ReadFromUDP(buf[:])
	
	if (addr != "" && useraddr != clientaddr) {
		return
	}//访问者与绑定者地址不同，不处理请求
	var taraddr string
	var portpos int
	if (buf[2] != 0x00) {
		fmt.Println("Fragment method not implemented")
		return
	}
	switch buf[3] {
		case 0x01:	{//IPV4
			addr = fmt.Sprintf("%d.%d.%d.%d", int(buf[4]), int(buf[5]), int(buf[6]), int(buf[7]))
			portpos = 8
		}
		case 0x03: {//DOMAIN NAME
			len := int(buf[4])
			for i := 0; i < len; i++ {
				addr += string(buf[i + 5])
			}
			portpos = len + 5
		}
		case 0x04:{//IPV6
			addr = fmt.Sprintf("%02x:%2x:%02x:%02x:%02x:%02x:%02x:%02x", 
				int(buf[4]), int(buf[6]), int(buf[8]), int(buf[10]), int(buf[12]), int(buf[14]), int(buf[16]), int(buf[18]))
			portpos = 20
		}
	}
	port := binary.BigEndian.Uint16(buf[portpos:portpos + 2])
	taraddr = fmt.Sprintf("%s:%d", taraddr, port)

	UDPServerAddr, err := net.ResolveUDPAddr("udp", taraddr)
	serverconn, err := net.ListenUDP("udp", UDPServerAddr)
	serverconn.WriteToUDP(buf[:], UDPServerAddr)
	count, _, err = serverconn.ReadFromUDP(buf[:])
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	clientconn.WriteToUDP(buf[:count], clientaddr)
	defer clientconn.Close()
	defer serverconn.Close()
}