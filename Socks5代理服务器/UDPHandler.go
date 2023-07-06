package Socks5

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
	switch buf[4] {
		case 0x01:	{//IPV4
			taraddr = fmt.Sprintf("%d.%d.%d.%d", int(buf[5]), int(buf[6]), int(buf[7]), int(buf[8]))
			portpos = 9
		}
		case 0x03: {//DOMAIN NAME
			len := int(buf[5])
			for i := 0; i < len; i++ {
				taraddr += string(buf[i + 6])
			}
			portpos = len + 6
		}
		case 0x04:{//IPV6
			taraddr = fmt.Sprintf("%02x:%2x:%02x:%02x:%02x:%02x:%02x:%02x", 
				int(buf[5]), int(buf[7]), int(buf[9]), int(buf[11]), int(buf[13]), int(buf[15]), int(buf[17]), int(buf[19]))
			portpos = 21
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
}