package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func Is_HTTP_Request(buf []byte, count int) int {
	if count > 4 {
		count = 4
	}
	str := string(buf[:count])
	if strings.Contains(str, "GET") {
		return 0
	} else if strings.Contains(str, "POST") {
		return 1
	} else {
		return 2
	}
}

func Forward(tarConn, srcConn net.Conn) {
	buf := make([]byte, 32*1024)
	HTTPtag := -1
	for {
		rcount, err := srcConn.Read(buf[:32*1024])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if HTTPtag == -1 {
			HTTPtag = Is_HTTP_Request(buf, rcount)
			if HTTPtag == 0 {
				fmt.Println("HTTP GET")
			} else if HTTPtag == 1 {
				fmt.Println("HTTP POST")
			} else {
				fmt.Println("NOT HTTP")
			}
		}
		if rcount > 0 {
			wcount, err := tarConn.Write(buf[:rcount])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if wcount < rcount {
				fmt.Println("End of forward")
				return
			}
		}
	}
}

func HandleTLSConnect(tar string, clientConn net.Conn) {
	defer clientConn.Close()
	// if err != nil {
	// 	log.Printf("proxy: dial: %s", err)
	// 	return
	// }
	// defer serverConn.Close()

	// config := &tls.Config{
	// 	MinVersion: tls.VersionTLS11,
	// 	MaxVersion: tls.VersionTLS13,
	// }

	serverConn, err := net.Dial("tcp", tar)
	//这里不能使用tls.Dial，因为tls.Dial返回的是对应tls的连接

	if err != nil {
		log.Printf("proxy: dial: %s", err)
		return
	}
	defer serverConn.Close()
	clientConn.Write(([]byte)("HTTP/1.1 200 xyzzy\r\nContent-Length: 0\r\n\r\n"))

	go io.Copy(clientConn, serverConn)
	io.Copy(serverConn, clientConn)
}

func HandleConnect(tar string, clientConn net.Conn) {
	fmt.Println("CONNECT success!!!")
	tarconn, err := net.DialTimeout("tcp", tar, 3*time.Second)
	defer clientConn.Close()
	defer tarconn.Close()
	//与目标连接
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			AnswerRequest(clientConn, 5)
		}
		if strings.Contains(err.Error(), "lookup") {
			AnswerRequest(clientConn, 4)
		} else if strings.Contains(err.Error(), "network is unreachable") {
			AnswerRequest(clientConn, 3)
		}
		return
	}
	AnswerRequest(clientConn, 0)
	//回复代理请求
	// var buf [512]byte
	// conn.Read(buf[:])
	go io.Copy(clientConn, tarconn)
	Forward(tarconn, clientConn)
	// go io.Copy(conn, tarconn)
	// io.Copy(tarconn, conn)
}
