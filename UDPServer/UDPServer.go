package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:	net.IPv4(0, 0, 0, 0),
		Port: 1234,
	})//本UDP服务器监听1234端口
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	var buf[512] byte
	var str = "The message accepted at UDP server is "
	for {
		count, addr, err := listener.ReadFromUDP(buf[:])
		if (err != nil) {
			fmt.Println(err.Error())
			return
		}
		sendstr := str + string(buf[:count])
		_, err = listener.WriteToUDP([]byte(sendstr), addr)
		fmt.Println(buf[:count])
		fmt.Println("Request coming from" + addr.String())
		if (err != nil) {
			fmt.Println(err.Error())
			return
		}
	}
}