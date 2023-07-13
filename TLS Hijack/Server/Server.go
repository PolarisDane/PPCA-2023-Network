package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.InsecureSkipVerify = true
	listener, err := tls.Listen("tcp", "localhost:1234", &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		go handleClient(conn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()
	var buf [512]byte
	count, err := clientConn.Read(buf[:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(buf[:count])
	clientConn.Write([]byte("Hello"))
}
