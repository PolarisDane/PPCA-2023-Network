package main

import (
	"fmt"
	"net"
	"io"
	"strings"
	"time"
)

func HandleConnect(tar string, conn net.Conn) {
	fmt.Println("CONNECT success!!!")
	tarconn, err := net.DialTimeout("tcp", tar, 3 * time.Second)
	//与目标连接
	if (err != nil){
		if (strings.Contains(err.Error(), "connection refused")) {
			AnswerRequest(conn, 5)
		}
		if (strings.Contains(err.Error(), "lookup")) {
			AnswerRequest(conn, 4)
		}else if (strings.Contains(err.Error(), "network is unreachable")) {
			AnswerRequest(conn, 3)
		}
		conn.Close()
		return
	}
	AnswerRequest(conn, 0)
	//回复代理请求
	go io.Copy(conn, tarconn)
	io.Copy(tarconn, conn)
	defer tarconn.Close()
}