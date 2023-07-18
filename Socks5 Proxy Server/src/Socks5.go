package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

var TLS_Hijack = true

func NegotiateAuthentication(conn net.Conn) (err error) {
	var buf [512]byte
	conn.Read(buf[:2])
	if int(buf[0]) != 0x05 {
		err = errors.New("Protocol version failed to match")
		return
	}
	NMETHODS := int(buf[1])
	conn.Read(buf[:NMETHODS])
	for i := 0; i < NMETHODS; i++ {
		if int(buf[i]) == 0x00 {
			return
		}
	}
	err = errors.New("Method not implemented")
	return
}

func AcceptRequest(conn net.Conn) (addr string, err error, CMD int) {
	var buf [512]byte
	_, err = io.ReadFull(conn, buf[:1])
	if err != nil {
		return
	}
	if buf[0] != 0x05 {
		err = errors.New("Protocol version failed to match")
		return
	}
	_, err = io.ReadFull(conn, buf[:1])
	if err != nil {
		return
	}
	switch buf[0] {
	case 0x01:
		{ //CONNECT
			CMD = 1
		}
	case 0x02:
		{ //BIND
			err = errors.New("CMD not supported")
			return
		}
	case 0x03:
		{ //UDP ASSOCIATE
			CMD = 3
		}
	default:
		{
			err = errors.New("CMD not supported")
			return
		}
	}
	_, err = io.ReadFull(conn, buf[:1])
	if err != nil {
		return
	}
	_, err = io.ReadFull(conn, buf[:1])
	if err != nil {
		return
	}
	switch buf[0] {
	case 0x01:
		{ //IPV4
			_, err = io.ReadFull(conn, buf[:4])
			if err != nil {
				return
			}
			addr = fmt.Sprintf("%d.%d.%d.%d", int(buf[0]), int(buf[1]), int(buf[2]), int(buf[3]))
		}
	case 0x03:
		{ //DOMAIN NAME
			_, err = io.ReadFull(conn, buf[:1])
			if err != nil {
				return
			}
			len := int(buf[0])
			_, err = io.ReadFull(conn, buf[:len])
			if err != nil {
				return
			}
			for i := 0; i < len; i++ {
				addr += string(buf[i])
			}
		}
	case 0x04:
		{ //IPV6
			_, err = io.ReadFull(conn, buf[:16])
			if err != nil {
				return
			}
			addr = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]",
				int(buf[0]), int(buf[1]), int(buf[2]), int(buf[3]), int(buf[4]), int(buf[5]), int(buf[6]), int(buf[7]),
				int(buf[8]), int(buf[9]), int(buf[10]), int(buf[11]), int(buf[12]), int(buf[13]), int(buf[14]), int(buf[15]))
		}
	default:
		{
			err = errors.New("ATYP not supported")
			return
		}
	}
	_, err = io.ReadFull(conn, buf[:2])
	if err != nil {
		return
	}
	port := binary.BigEndian.Uint16(buf[:2])
	addr = fmt.Sprintf("%s:%d", addr, port)
	return
}

func AnswerRequest(conn net.Conn, ErrorType int) {
	var buf [512]byte
	buf[0] = byte(0x05)
	switch ErrorType {
	case 0:
		buf[1] = byte(0x00)
	case 3:
		buf[1] = byte(0x03)
	case 4:
		buf[1] = byte(0x04)
	case 5:
		buf[1] = byte(0x05)
	case 7:
		buf[1] = byte(0x07)
	case 8:
		buf[1] = byte(0x08)
	}
	buf[2] = byte(0x00)
	buf[3] = byte(0x01)
	buf[4] = byte(0x00)
	buf[5] = byte(0x00)
	buf[6] = byte(0x00)
	buf[7] = byte(0x00)
	buf[8] = byte(0x00)
	buf[9] = byte(0x00)
	conn.Write(buf[:10])
}

func HandleConn(conn net.Conn) {
	var buf [4096]byte
	err := NegotiateAuthentication(conn)
	//协商认证
	if err != nil {
		buf[0] = 0x05
		buf[1] = 0xff
		conn.Write(buf[:2])
		//协商认证失败
		return
	}

	buf[0] = 0x05
	buf[1] = 0x00
	conn.Write(buf[:2])
	//完成协商认证
	fmt.Println("Accepting")
	tar, err, CMD := AcceptRequest(conn)
	//接受代理请求
	if err != nil {
		if err.Error() == "CMD not supported" {
			AnswerRequest(conn, 7)
		}
		if err.Error() == "ATYP not supported" {
			AnswerRequest(conn, 8)
		}
		conn.Close()
		return
	}
	fmt.Println("Accepted")
	if CMD == 1 {
		HandleConnect(tar, conn)
	} else if CMD == 3 {
		HandleUDP(tar, conn)
	}
}

func TCPLink(addr string) error {
	fmt.Println("Listening here")
	//非TLS劫持，不能处理HTTPS协议
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		fmt.Println("Proxy user is found")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		go HandleConn(conn)
	}
}

func main() {
	err := TCPLink(":8080")
	fmt.Println(err.Error())
}
