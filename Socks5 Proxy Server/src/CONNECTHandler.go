package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func Forward(tarConn, srcConn net.Conn) {
	buf := make([]byte, 32*1024)
	for {
		rcount, err := srcConn.Read(buf[:32*1024])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if rcount > 0 {
			wcount, err := tarConn.Write(buf[:rcount])
			//fmt.Println(string(buf[:rcount]))
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

func HandleConnect(tar string, clientConn net.Conn) {
	fmt.Println("CONNECT success!!!")
	tarConn, err := net.DialTimeout("tcp", tar, 3*time.Second)
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
	if TLS_Hijack {
		reader := bufio.NewReader(clientConn)
		line, err := reader.Peek(8)
		if err != nil {
			fmt.Println("Internal server error" + err.Error())
			return
		}
		if (line[0] == 0x16) && (line[1] == 0x03) && (line[2] == 0x01) {
			defer clientConn.Close()
			tarConn.Close()
			Listener, err := net.Listen("tcp", ":0")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			handshakeAddr := Listener.Addr().String()
			fmt.Println("Listening for TLS handshake")
			go func() {
				handshakeConn, err := Listener.Accept()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				cert, err := generateCert(tar)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				tlsConfig := &tls.Config{
					Certificates:       []tls.Certificate{cert},
					InsecureSkipVerify: true,
				}
				tlsConn := tls.Server(handshakeConn, tlsConfig)
				err = tlsConn.Handshake()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Received TLS handshake")

				tlsConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
				serverConn, err := tls.Dial("tcp", tar, tlsConfig)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				defer serverConn.Close()
				go ModifyForward(tlsConn, serverConn)
				io.Copy(serverConn, tlsConn)
			}()
			proxyConn, err := net.DialTimeout("tcp", handshakeAddr, 3*time.Second)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer proxyConn.Close()
			var buf [32 * 1024]byte
			go io.Copy(clientConn, proxyConn)
			n, err := reader.Read(buf[:])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			proxyConn.Write(buf[:n])
			io.Copy(proxyConn, clientConn)
		}
	} else {
		defer clientConn.Close()
		defer tarConn.Close()
		go ModifyForward(clientConn, tarConn)
		io.Copy(tarConn, clientConn)
	}
}
