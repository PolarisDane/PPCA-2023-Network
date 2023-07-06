package Socks5

import (
	"fmt"
	"net"
	"errors"
	"encoding/binary"
)

func NegotiateAuthentication (conn net.Conn) (result int, error string) {
	var buf[512] byte
	conn.Read(buf[:])
	if (int(buf[0]) != 0x05) {
		error = "Protocol version failed to match"
		return
	}
	for i := 0; i < int(buf[1]); i++ {
		if (int(buf[i + 1]) == 0x00) {
			return
		}
	}
	error = "Method not implemented"
	return
}

func AcceptRequest(conn net.Conn) (addr string, err error, RequestType int) {
	var buf [512] byte
	conn.Read(buf[:])
	fmt.Println(buf[0]);
	if (int(buf[0]) != 0x05) {
		err = errors.New("Protocol version failed to match")
		return
	}
	switch buf[1] {
		case 0x01: {//CONNECT
			RequestType = 1
		}
		case 0x02: {//BIND
			RequestType = 2
			err = errors.New("Proxy method not implemented")
			return
		}
		case 0x03: {//UDP ASSOCIATE
			RequestType = 3
		}
	}
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
			addr = fmt.Sprintf("%02x:%2x:%02x:%02x:%02x:%02x:%02x:%02x", 
				int(buf[4]), int(buf[6]), int(buf[8]), int(buf[10]), int(buf[12]), int(buf[14]), int(buf[16]), int(buf[18]))
			portpos = 20
		}
	}
	port := binary.BigEndian.Uint16(buf[portpos:])
	addr = fmt.Sprintf("%s:%d", addr, port)
	return
}

func HandleConn(conn net.Conn) {
	var buf [4096]byte
	fmt.Println("Authenticating")
	NegotiateAuthentication(conn)
	//协商认证
	buf[0] = 0x05
	buf[1] = 0x00
	conn.Write(buf[:2])
	//完成协商认证
	fmt.Println("Accepting request")
	tar, err, RequestType := AcceptRequest(conn)
	fmt.Println("Accepted")
	//接受代理请求
	if (err != nil) {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	if (RequestType == 1) {
		HandleConnect(tar, conn)
	} else if (RequestType == 3) {
		fmt.Println("UDP success!!!")
		HandleUDP(tar, conn)
		//对于UDP ASSOCIATE来说请求代理中的地址和端口指的是客户端发送UDP包的地址和端口
	}
	defer conn.Close()
}

func TCPLink(addr string) error {
	fmt.Println("Listening here")
	listener, err := net.Listen("tcp", addr)
	if (err != nil) {
		return err
	}
	for {
		conn, err := listener.Accept()
		fmt.Println("Proxy user is found")
		if (err != nil) {
			return err
		}
		go HandleConn(conn)
	}
}

func main() {
	err := TCPLink(":8080")
	fmt.Println(err.Error())
}
