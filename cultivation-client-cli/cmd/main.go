package main

import (
	"fmt"
	"strings"

	"cultivation-client-cli/internal/client"
	"cultivation-client-cli/internal/commands"

	"github.com/chzyer/readline"
)

func main() {
	// Single readline for both auth and game (avoids stdin conflicts)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "选择 [1]登录 [2]注册 [3]退出: ",
		HistoryFile: ".cli_history",
	})
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}
	defer rl.Close()

	var token, entityID string
authLoop:
	for {
		line, err := rl.Readline()
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		switch line {
		case "1", "login":
			rl.SetPrompt("用户名: ")
			user, _ := rl.Readline()
			user = strings.TrimSpace(user)
			rl.SetPrompt("密码: ")
			pass, _ := rl.Readline()
			pass = strings.TrimSpace(pass)

			var err error
			token, entityID, err = client.Login(user, pass)
			if err != nil {
				fmt.Printf("登录失败: %v\n", err)
				rl.SetPrompt("选择 [1]登录 [2]注册 [3]退出: ")
				continue authLoop
			}
			break authLoop

		case "2", "register":
			rl.SetPrompt("用户名: ")
			user, _ := rl.Readline()
			user = strings.TrimSpace(user)
			rl.SetPrompt("密码: ")
			pass, _ := rl.Readline()
			pass = strings.TrimSpace(pass)

			var err error
			token, entityID, err = client.Register(user, pass)
			if err != nil {
				fmt.Printf("注册失败: %v\n", err)
				rl.SetPrompt("选择 [1]登录 [2]注册 [3]退出: ")
				continue authLoop
			}
			break authLoop

		case "3", "exit", "quit":
			return
		}
	}

	msgCh := make(chan string, 64)
	doneCh := make(chan struct{})

	conn, err := client.ConnectWebSocket(token, msgCh, doneCh)
	if err != nil {
		fmt.Printf("WebSocket连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	// Wait for initial state_sync
	<-msgCh

	// Reuse readline for game loop with screen manager
	client.SetReadline(rl)
	client.SetCompletionFunc(commands.GetCompletions)
	screen, err := client.NewScreen(msgCh, doneCh)
	if err != nil {
		fmt.Printf("界面初始化失败: %v\n", err)
		return
	}
	defer screen.Close()

	screen.Start(func(line string) {
		commands.Dispatch(conn, entityID, line)
	})
}
