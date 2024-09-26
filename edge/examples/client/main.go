package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"gozeroImSystem/common/libnet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:15678") // 之前的443端口，服务端那边收不到
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Enter text to send:")

	protocol := libnet.NewIMProtocol()
	codec := protocol.NewCodec(conn)

	go readServerResponse(codec)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("client: ")
		text, _ := reader.ReadString('\n')
		msg := libnet.Message{
			Body: []byte(text),
		}
		err = codec.Send(msg)
		if err != nil {
			fmt.Printf("send error: %v\n", err)
		}
		time.Sleep(1000 * time.Microsecond)
	}
}

func readServerResponse(codec libnet.Codec) {
	var tmpDelay time.Duration
	for {
		msg, err := codec.Receive()
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from server:", err)
			break
		}
		if err == io.EOF {
			if tmpDelay == 0 {
				// 第一次等待重试时长为10ms
				tmpDelay = 10 * time.Microsecond
			} else {
				// 之后每次等待重试时长x2
				tmpDelay *= 2
			}
			// 等待时长最多为1s
			if max := time.Second; tmpDelay > max {
				tmpDelay = max
			}
			time.Sleep(tmpDelay)
			continue
		}
		fmt.Println("Server response: " + string(msg.Body))
	}
}
