package Socks5

import (
	"fmt"
	"net"
	"io"
	"time"
)

func AnswerConnectRequest(conn net.Conn) {
	var buf[512] byte
	buf[0] = byte(0x05)
	buf[1] = byte(0x00)
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

func HandleConnect(tar string, conn net.Conn) {
	fmt.Println("CONNECT success!!!")
	tarconn, err := net.DialTimeout("tcp", tar, 3 * time.Second)
	//与目标连接
	if (err != nil){
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	AnswerConnectRequest(conn)
	//回复代理请求
	go io.Copy(tarconn, conn)
	io.Copy(conn, tarconn)
	defer tarconn.Close()
}