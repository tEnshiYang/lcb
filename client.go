package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("dial failed, err :", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console,err: %v\n", err)
			break
		}
		data = strings.TrimSpace(data)
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Printf("write failed ,err:%v\n", err)
			break
		}
		buf := [128]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
