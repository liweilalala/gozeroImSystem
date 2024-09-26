package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 1. 监听端口
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// 2. 接受连接
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		// 3. 启动一个 goroutine 处理连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 启动一个 goroutine 发送消息给客户端
	go func() {
		for {
			// 从标准输入读取服务器端输入的消息
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Server: ")
			message, _ := reader.ReadString('\n')

			// 发送消息到客户端
			_, err := conn.Write([]byte(message))
			if err != nil {
				log.Println("Failed to write to connection:", err)
				return
			}
		}
	}()

	// 接收客户端发送的消息
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Connection closed by client:", err)
			return
		}

		fmt.Print("Client: " + string(buffer[:n]))
	}
}
