package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("listen failed,err:", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed,err:%v\n", err)
			continue
		}
		go process(conn)
	}

}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		str := string(buf[:n])
		fmt.Printf("recv from client, data: %v\n", str)
		res := add(str)
		conn.Write([]byte(res))
	}
}

func add(str string) string {
	idx := strings.Index(str, "+")
	if idx == -1 {
		return "0"
	}
	a, err := strconv.Atoi(str[:idx])
	if err != nil {
		return "0"
	}
	b, err := strconv.Atoi(str[idx:])
	if err != nil {
		return "0"
	}

	return strconv.Itoa(a + b)
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &x)
	return bytesBuffer.Bytes()
}
