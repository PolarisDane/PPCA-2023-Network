package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func ModifyForward(tarConn, srcConn net.Conn) {
	data := make([]byte, 0)
	buf := make([]byte, 32*1024)
	tag := -1
	for {
		rcount, err := srcConn.Read(buf[:32*1024])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if tag == -1 {

		}
	}
}

func Is_HTTP_Request(buf []byte, count int) int {
	if count > 8 {
		count = 8
	}
	str := string(buf[:count])
	switch {
	case strings.Contains(str, "GET"):
		return 1
	case strings.Contains(str, "POST"):
		return 2
	case strings.Contains(str, "PUT"):
		return 3
	case strings.Contains(str, "DELETE"):
		return 4
	case strings.Contains(str, "PATCH"):
		return 5
	case strings.Contains(str, "HEAD"):
		return 6
	case strings.Contains(str, "OPTIONS"):
		return 7
	default:
		return 0
	}
}

func Http_Response_Parse(buf []byte, n int) int64 {
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(buf)), nil)
	if err != nil {
		return -1
	}
	return resp.ContentLength
}

func Http_Response_Modify(buf []byte, n int) ([]byte, int) {
	index := strings.Index(string(buf), "\r\n\r\n")
	fmt.Printf("Origin2:\n%v\n", string(buf))
	data := make([]byte, 0)
	data = append(data, buf[index:]...)

	// Here, do something to data.
	// data = bytes.Replace(data, []byte("努力"), []byte("天天摆烂"), -1)

	head := Http_Response_Modify_Head(string(buf[:index]), len(data)-4)

	head = append(head, data...)
	fmt.Printf("Final2:\n%v\n", string(head))
	fmt.Println("Compare", len(head)-len(buf))

	return head, len(head)
}

func Http_Response_Modify_Head(str string, n int) []byte {
	slice := strings.Split(str, "\r\n")
	target := "Content-Length: " + fmt.Sprintf("%d", n)
	for i := 0; i < len(slice); i++ {
		if strings.HasPrefix(slice[i], "Content-Length:") {
			fmt.Println("Code:", slice[i])
			slice[i] = target
			break
		}
		if strings.HasPrefix(slice[i], "Transfer-Encoding: chunked") {
			fmt.Println("Code:", slice[i])
			slice[i] = target
			break
		}
	}
	return []byte(strings.Join(slice, "\r\n"))
}
