package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cultivation-bubbletea/internal/client"
	"cultivation-bubbletea/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("=== 修仙世界 MUD - Bubble Tea TUI ===")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var token, entityID string

authLoop:
	for {
		fmt.Print("选择 [1]登录 [2]注册 [3]退出: ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		var username, password string

		switch choice {
		case "1", "login":
			fmt.Print("用户名: ")
			if !scanner.Scan() {
				return
			}
			username = strings.TrimSpace(scanner.Text())

			fmt.Print("密码: ")
			if !scanner.Scan() {
				return
			}
			password = strings.TrimSpace(scanner.Text())

			var err error
			token, entityID, err = client.Login(username, password)
			if err != nil {
				fmt.Printf("登录失败: %v\n", err)
				continue authLoop
			}
			break authLoop

		case "2", "register":
			fmt.Print("用户名: ")
			if !scanner.Scan() {
				return
			}
			username = strings.TrimSpace(scanner.Text())

			fmt.Print("密码: ")
			if !scanner.Scan() {
				return
			}
			password = strings.TrimSpace(scanner.Text())

			var err error
			token, entityID, err = client.Register(username, password)
			if err != nil {
				fmt.Printf("注册失败: %v\n", err)
				continue authLoop
			}
			break authLoop

		case "3", "exit", "quit":
			fmt.Println("再见!")
			return
		}
	}

	// Connect WebSocket
	fmt.Print("正在连接服务器...")

	msgCh := make(chan client.WsMessage, 128)
	state := &client.GameState{}

	conn, err := client.ConnectWebSocket(token, msgCh, state)
	if err != nil {
		fmt.Printf("\nWebSocket连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	// Wait for the initial state_sync message
	<-msgCh

	fmt.Println(" 已连接!")
	fmt.Println("按 Ctrl+C 退出游戏")
	fmt.Println()

	// Create and start Bubble Tea program
	m := tea.NewProgram(model.NewWithChannel(conn, entityID, msgCh, state),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := m.Run(); err != nil {
		fmt.Printf("程序错误: %v\n", err)
		os.Exit(1)
	}
}
