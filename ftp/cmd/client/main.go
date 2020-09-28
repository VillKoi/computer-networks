package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/VillKoi/computer-networks/ftp/client"
	"golang.org/x/sync/errgroup"
)

const (
	login    = "vill"
	password = "qwerty"

	httpport = ":9000"
)

func InvalidCommand() {
	fmt.Println("invalid command")
}

func main() {

	// c := new(client.FTPClient)
	// var err error

	// c.FTP, err = ftp.Dial("localhost:2121", ftp.DialWithTimeout(5*time.Second))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = c.FTP.Login(login, password)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer c.FTP.Quit()

	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	fmt.Println("Enter command: ")
	// 	text, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("GetLines: " + err.Error())
	// 		break
	// 	}

	// 	cmd := strings.Split(strings.TrimSpace(text), " ")
	// 	if len(cmd) == 0 {
	// 		continue
	// 	}

	// 	switch cmd[0] {
	// 	case "put":
	// 		if len(cmd) != 3 {
	// 			InvalidCommand()
	// 			continue
	// 		}
	// 		err = c.StoreFile(cmd[1], cmd[2])
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "get":
	// 		if len(cmd) != 3 {
	// 			InvalidCommand()
	// 			continue
	// 		}
	// 		err = c.ReadFile(cmd[1], cmd[2])
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "pwd":
	// 		err = c.CurrentDir()
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "ls":
	// 		var path string
	// 		if len(cmd) == 2 {
	// 			path = cmd[1]
	// 		}
	// 		err = c.List(path)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "cd":
	// 		if len(cmd) != 2 {
	// 			InvalidCommand()
	// 			continue
	// 		}
	// 		err := c.ChangeDir(cmd[1])
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "mkdir":
	// 		if len(cmd) != 2 {
	// 			InvalidCommand()
	// 			continue
	// 		}
	// 		err := c.MakeDir(cmd[1])
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	case "delete":
	// 		if len(cmd) != 2 {
	// 			InvalidCommand()
	// 			continue
	// 		}
	// 		err := c.Delete(cmd[1])
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	default:
	// 		InvalidCommand()
	// 	}
	// }

	ctx, cancel := context.WithCancel(context.Background())

	c := new(client.FTPClient)
	var err error

	group := errgroup.Group{}
	group.Go(func() error {
		return StartHTTP(ctx, c)
	})

	group.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
		return nil
	})

	err = group.Wait()
	if err != nil {
		log.Fatal("error from group: ", err)
	}
}
