package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
	if err != nil {
		log.Fatalf("proxy: loadkeys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.InsecureSkipVerify = true
	listener, err := tls.Listen("tcp", "localhost:8080", &config)
	if err != nil {
		log.Fatalf("proxy: listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("proxy: accept: %s", err)
			break
		}
		go handleClient(conn)
	}
}

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
			fmt.Print(string(buf[:rcount]))
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

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()
	// if err != nil {
	// 	log.Printf("proxy: dial: %s", err)
	// 	return
	// }
	// defer serverConn.Close()

	var buf [512]byte

	n, err := clientConn.Read(buf[:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var str = string(buf[:n])

	fmt.Println(string(buf[:n]))

	var splitstr = strings.Split(str, " ")

	var addr = splitstr[1]

	fmt.Println(addr)

	// config := &tls.Config{
	// 	MinVersion: tls.VersionTLS11,
	// 	MaxVersion: tls.VersionTLS13,
	// }

	serverConn, err := net.Dial("tcp", addr)
	//这里不能使用tls.Dial，因为tls.Dial返回的是对应tls的连接

	if err != nil {
		log.Printf("proxy: dial: %s", err)
		return
	}
	defer serverConn.Close()
	clientConn.Write(([]byte)("HTTP/1.1 200 xyzzy\r\nContent-Length: 0\r\n\r\n"))

	go Forward(clientConn, serverConn)
	io.Copy(serverConn, clientConn)

}
