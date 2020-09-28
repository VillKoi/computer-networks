package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vill_koi/client"

	"github.com/jlaffaye/ftp"
)

const (
	login    = "vill"
	password = "qwerty"
)

func InvalidCommand() {
	fmt.Println("invalid command")
}

func main() {
	client := new(client.FTPClient)
	var err error

	client.ftp, err = ftp.Dial("localhost:2121", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = client.ftp.Login(login, password)
	if err != nil {
		log.Fatal(err)
	}
	defer client.ftp.Quit()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter command: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("GetLines: " + err.Error())
			break
		}

		cmd := strings.Split(strings.TrimSpace(text), " ")
		if len(cmd) == 0 {
			continue
		}

		switch cmd[0] {
		case "put":
			if len(cmd) != 3 {
				InvalidCommand()
				continue
			}
			err = client.StoreFile(cmd[1], cmd[2])
			if err != nil {
				fmt.Println(err)
			}
		case "get":
			if len(cmd) != 3 {
				InvalidCommand()
				continue
			}
			err = client.ReadFile(cmd[1], cmd[2])
			if err != nil {
				fmt.Println(err)
			}
		case "pwd":
			err = client.CurrentDir()
			if err != nil {
				fmt.Println(err)
			}
		case "ls":
			var path string
			if len(cmd) == 2 {
				path = cmd[1]
			}
			err = client.List(path)
			if err != nil {
				fmt.Println(err)
			}
		case "cd":
			if len(cmd) != 2 {
				InvalidCommand()
				continue
			}
			err := client.ChangeDir(cmd[1])
			if err != nil {
				fmt.Println(err)
			}
		case "mkdir":
			if len(cmd) != 2 {
				InvalidCommand()
				continue
			}
			err := client.MakeDir(cmd[1])
			if err != nil {
				fmt.Println(err)
			}
		default:
			InvalidCommand()
		}
	}
}
