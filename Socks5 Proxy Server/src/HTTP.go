package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func ForwardToClient(tarConn, srcConn net.Conn) {
	data := make([]byte, 0)
	buf := make([]byte, 32*1024)
	tag := -1
	lens := 0
	index := -1
	contentLength := int64(-1)
	compress := 0
	for {
		rcount, err := srcConn.Read(buf[:32*1024])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if tag == -1 {
			tag = IsResponse(buf, rcount)
		}
		if tag > 0 {
			lens += rcount
			data = append(data, buf[:rcount]...)
			if index == -1 {
				index = strings.Index(string(data), "\r\n\r\n")
				//找到响应头
				if index != -1 {
					contentLength = ResponseParse(data)
					if contentLength == -2 {
						return
					}
				}
			}
			if index != -1 {
				if contentLength == -1 {
					//内容长度不确定，如使用chunked
					if strings.HasSuffix(string(data[lens-5:lens]), "0\r\n\r\n") {
						//chunked传输的结尾为0\r\n\r\n
						data, lens, compress = ResponseDecode(data, lens)
						if lens == -1 {
							return
						}
						tag = -1
					}
				} else {
					if lens == index+int(contentLength)+4 {
						tag = -1
					}
				}
			}
			if tag == -1 {
				data, lens = ResponseModify(data, lens, compress)
				wcount, err := tarConn.Write(data[:lens])
				//fmt.Println(string(data[:lens]))
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				if wcount < lens {
					fmt.Println("Package end")
					return
				}
				data = data[:0]
				lens = 0
				index = -1
				contentLength = -1
				compress = 0
			}
		} else {
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
}

func IsResponse(buf []byte, count int) int {
	if count > 8 {
		count = 8
	}
	if strings.HasPrefix(string(buf[:count]), "HTTP") {
		return 1
	} else {
		return 0
	}
}

func IsRequest(buf []byte, count int) int {
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

func ResponseParse(buf []byte) int64 {
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(buf)), nil)
	if err != nil {
		fmt.Println(err.Error())
		return -2
	}
	return resp.ContentLength
}

func ResponseDecode(buf []byte, count int) ([]byte, int, int) {
	str := string(buf)
	index := strings.Index(str, "\r\n\r\n") + 4
	data := make([]byte, 0)
	head := make([]byte, 0)
	ind := strings.Index(str[:index], "gzip")
	compress := 0
	if ind >= 0 {
		compress = 1
	}
	ind = strings.Index(str[:index], "br")
	if ind >= 0 {
		compress = 2
	}
	ind = strings.Index(str[:index], "deflate")
	if ind >= 0 {
		compress = 3
	}
	head = append(head, str[:index]...)
	for {
		siz := 0
		fmt.Sscanf(str[index:], "%x", &siz)
		if siz == 0 {
			break
		}
		index += strings.Index(str[index:], "\r\n") + 2
		data = append(data, str[index:index+siz]...)
		index += siz + 2
	}
	head = append(head, data...)
	return head, 0, compress
} //将chunked方式传输的数据拼接

func ResponseModify(buf []byte, n int, compress int) ([]byte, int) {
	index := strings.Index(string(buf), "\r\n\r\n") + 4
	//fmt.Printf("Origin2:\n%v\n", string(buf))
	data := make([]byte, 0)
	data = append(data, buf[index:]...)

	body := ResponseBodyModify(buf[index:], compress)

	head := ResponseHeadModify(string(buf[:index]), len(data)-4)

	head = append(head, body...)
	//fmt.Printf("Final2:\n%v\n", string(head))
	//fmt.Println("Compare", len(head)-len(buf))

	return head, len(head)
}

func ResponseBodyModify(data []byte, compress int) []byte {
	if len(data) == 0 {
		return data
	}
	ret := make([]byte, 0)
	ret = append(ret, data...)
	switch compress {
	case 1:
		ret = Decompressgzip(ret)
	case 2:
		ret = Decompressbr(ret)
	case 3:
		ret = Decompressflate(ret)
	}
	ret = bytes.Replace(ret, []byte("百度"), []byte("Polaris_Dane"), -1)
	switch compress {
	case 1:
		ret = Compressgzip(ret)
	case 2:
		ret = Compressbr(ret)
	case 3:
		ret = Compressflate(ret)
	}
	return ret
}

func ResponseHeadModify(str string, n int) []byte {
	slice := strings.Split(str, "\r\n")
	target := "Content-Length: " + fmt.Sprintf("%d", n)
	for i := 0; i < len(slice); i++ {
		if strings.HasPrefix(slice[i], "Content-Length:") {
			slice[i] = target
			break
		}
		if strings.HasPrefix(slice[i], "Transfer-Encoding: chunked") {
			slice[i] = target
			break
		}
	}
	return []byte(strings.Join(slice, "\r\n"))
}
