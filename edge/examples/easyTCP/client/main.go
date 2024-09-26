package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 启动一个 goroutine 接收服务器发送的消息
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				log.Println("Connection closed by server:", err)
				return
			}

			fmt.Print("Server: " + string(buffer[:n]))
		}
	}()

	// 发送消息给服务器
	for {
		// 从标准输入读取客户端输入的消息
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Client: ")
		message, _ := reader.ReadString('\n')

		// 发送消息到服务器
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Failed to write to connection:", err)
			return
		}
	}
}
