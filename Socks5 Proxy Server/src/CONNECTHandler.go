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

func Is_HTTP_Request(buf []byte, count int) int {
	if count > 16 {
		count = 16
	}
	str := string(buf[:count])
	if strings.Contains(str, "GET") {
		return 0
	} else if strings.Contains(str, "POST") {
		return 1
	} else if strings.Contains(str, "CONNECT") {
		return 2
	} else {
		return 3
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
			} else if HTTPtag == 2 {
				fmt.Println("HTTP CONNECT")
			} else {
				fmt.Println("NOT HTTP")
			}
		}
		if rcount > 0 {
			wcount, err := tarConn.Write(buf[:rcount])
			fmt.Println(string(buf[:rcount]))
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
	// defer clientConn.Close()

	// serverConn, err := tls.Dial("tcp", tar, &tls.Config{
	// 	InsecureSkipVerify: true,
	// })
	// if err != nil {
	// 	log.Printf("proxy: dial: %s", err)
	// 	return
	// }

	// cert =

	// defer serverConn.Close()
	// clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// go Forward(clientConn, serverConn)
	// io.Copy(serverConn, clientConn)
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
	reader := bufio.NewReader(clientConn)
	line, err := reader.Peek(8)
	fmt.Println("Hello!!!")
	if err != nil {
		fmt.Println("Internal server error" + err.Error())
		return
	}
	if (line[0] == 0x16) && (line[1] == 0x03) && (line[2] == 0x01) && TLS_Hijack {
		fmt.Println("Bye!!!")
		tarconn.Close()
		Listener, err := net.Listen("tcp", "127.0.0.1:9000")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
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
			go Forward(tlsConn, serverConn)
			io.Copy(serverConn, tlsConn)
		}()
		// cert, err := generateCert(tar)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }
		// tlsConfig := &tls.Config{
		// 	Certificates:       []tls.Certificate{cert},
		// 	InsecureSkipVerify: true,
		// }
		// tlsClientConn := tls.Server(clientConn, tlsConfig)
		// defer tlsClientConn.Close()
		proxyConn, err := net.Dial("tcp", "127.0.0.1:9000")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer proxyConn.Close()
		var buf [32 * 1024]byte
		go Forward(clientConn, proxyConn)
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		proxyConn.Write(buf[:n])
		io.Copy(proxyConn, clientConn)
	} else {
		go io.Copy(clientConn, tarconn)
		Forward(tarconn, clientConn)
	}
}
