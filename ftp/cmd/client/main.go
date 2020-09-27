package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	login    = "vill"
	password = "qwerty"
)

type FTPClient struct {
	ftp *ftp.ServerConn
}

// put local/example.txt server/new_file.txt
func (c FTPClient) StoreFile(path, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)
	err = c.ftp.Stor(path, buffer)
	return err
}

// get example.txt new_file.txt
func (c FTPClient) ReadFile(path, newPath string) error {
	r, err := c.ftp.Retr(path)
	if err != nil {
		return err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(newPath, buf, 0644)
	return err
}

// pwd
func (c FTPClient) CurrentDir() error {
	pwd, err := c.ftp.CurrentDir()
	if err != nil {
		return err
	}
	fmt.Println(pwd)
	return nil
}

// ls
func (c FTPClient) List(path string) error {
	list, err := c.ftp.NameList(path)
	if err != nil {
		return err
	}
	fmt.Println(list)
	return nil
}

// cd
func (c FTPClient) ChangeDir(path string) error {
	err := c.ftp.ChangeDir(path)
	return err
}

// mkdir
func (c FTPClient) MakeDir(path string) error {
	err := c.ftp.MakeDir(path)
	return err
}

func InvalidCommand() {
	fmt.Println("invalid command")
}

func main() {
	client := new(FTPClient)
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
