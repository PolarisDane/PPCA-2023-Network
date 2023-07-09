package main

import (
	"fmt"
	"net"
	"encoding/binary"
	"strconv"
)

func AnswerUDPRequest(conn net.Conn, DistributedAddr net.Addr) {
	host, port, _ := net.SplitHostPort(DistributedAddr.String())
	DistributedIP := net.ParseIP(host)
	DistributedPort, _ := strconv.Atoi(port)
	fmt.Println(DistributedPort)
	var buf []byte
	buf = append(buf, byte(0x05))
	buf = append(buf, byte(0x00))
	buf = append(buf, byte(0x00))
	buf = append(buf, byte(0x04))
	buf = append(buf, DistributedIP...)
	fmt.Println(binary.BigEndian.AppendUint16(buf, uint16(DistributedPort)))
	conn.Write(binary.BigEndian.AppendUint16(buf, uint16(DistributedPort)))
}

func HandleUDP(addr string, conn net.Conn) {
	fmt.Println("UDP success!!!")
	var buf[512] byte
	clientaddr, _ := net.ResolveUDPAddr("udp", addr)
	clientconn, err := net.ListenUDP("udp", nil)
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	//系统随机分配端口
	var count int
	var useraddr *net.UDPAddr
	var sendport *net.UDPAddr
	var serverconn *net.UDPConn
	AnswerUDPRequest(conn, clientconn.LocalAddr())
	//回复代理请求
	go func() {
		var tmp[512] byte
		conn.Read(tmp[:])
		if (tmp[0] == 0x04) {
			clientconn.Close()
			serverconn.Close()
			return
		}
	}()
	defer clientconn.Close()
	for {
		//go func() {
			count, useraddr, err = clientconn.ReadFromUDP(buf[:])
			if (err != nil) {
				fmt.Println(err.Error())
				return
			}
			if (count == 1 && buf[0] == 0x04) {
				return
			}
			fmt.Println(buf[:count])
			fmt.Println(useraddr)
			//对于无连接的UDP，发送者可以直接使用Write，但是接受者回应时只能使用WriteToUDP
		
			if (addr != "" && useraddr.String() != clientaddr.String()) {
				return
			}//访问者与绑定者地址不同，不处理请求
			if (buf[2] != 0x00) {
				fmt.Println("Fragment method not implemented")
				return
			}
			var taraddr string
			var portpos int
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
					addr = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]", 
						int(buf[4]), int(buf[5]), int(buf[6]), int(buf[7]), int(buf[8]), int(buf[9]), int(buf[10]), int(buf[11]),
						int(buf[12]), int(buf[13]), int(buf[14]), int(buf[15]), int(buf[16]), int(buf[17]), int(buf[18]), int(buf[19]))
					portpos = 20
				}
			}
			port := binary.BigEndian.Uint16(buf[portpos:portpos + 2])
			taraddr = fmt.Sprintf("%s:%d", addr, port)
		//}()

		UDPServerAddr, err := net.ResolveUDPAddr("udp", taraddr)
		if (sendport == nil) {
			serverconn, err = net.ListenUDP("udp", nil)
			defer serverconn.Close()
			fmt.Println(serverconn.LocalAddr())
			sendport, err = net.ResolveUDPAddr("udp", serverconn.LocalAddr().String())
			//使用系统随机分配的端口，注意给记录下来保证对于一个客户端使用端口固定
		}
		
		serverconn.WriteToUDP(buf[portpos + 2:count], UDPServerAddr)
		fmt.Println(UDPServerAddr)
		//go func() {3
			count, _, err = serverconn.ReadFromUDP(buf[:])
		//}()
		if (err != nil) {
			fmt.Println(err.Error())
			return
		}
		clientconn.WriteToUDP(buf[:count], clientaddr)
	}
}