package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Server struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func sshservers(path, cmd string) bytes.Buffer {
	var servers []Server

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &servers)

	resChannel := make(chan string, len(servers))
	timeout := time.After(5 * time.Second)

	divider := fmt.Sprintf("cmd:%s", cmd)

	for i := 0; i < len(servers); i++ {
		go func(j int) {
			config := &ssh.ClientConfig{
				User: servers[j].User,
				Auth: []ssh.AuthMethod{
					ssh.Password(servers[j].Password)},
				HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
					return nil
				}}

			conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%d", servers[j].Host, servers[j].Port), config)
			session, _ := conn.NewSession()
			defer session.Close()

			stdin, _ := session.StdinPipe()
			stdout, _ := session.StdoutPipe()
			err = session.Shell()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintf(stdin, "echo %s\n", divider)
			fmt.Fprintf(stdin, "%s\n", cmd)
			fmt.Fprintf(stdin, "echo getpath\n")
			fmt.Fprintf(stdin, "pwd\n")
			fmt.Fprintf(stdin, "%s\n", "exit")

			out, _ := ioutil.ReadAll(stdout)
			outstr := string(out)
			outstr = outstr[strings.Index(outstr, divider+"\n")+len(divider)+1 : strings.Index(outstr, "getpath\n")]

			resChannel <- outstr
		}(i)
	}

	var stdout bytes.Buffer

	for i := 0; i < len(servers); i++ {
		select {
		case res := <-resChannel:
			stdout.WriteString(res)
		case <-timeout:
			stdout.WriteString("Timed out!\n")
		}
	}
	return stdout
}

func main() {
	var path string
	flag.StringVar(&path, "path", "package.json", "path to json file")

	var stdout bytes.Buffer

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		stdout = sshservers(path, cmd)

		if bytes.Compare([]byte(cmd), []byte("exit")) == 0 {
			break
		}
		fmt.Println(stdout.String())
		stdout.Reset()
	}
}