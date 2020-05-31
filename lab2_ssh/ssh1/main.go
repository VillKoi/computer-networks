package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

func main() {

	var (
		user     = flag.String("user", "iu9_32_06", "Username for login")
		password = flag.String("pass", "", "Password for login")
		host     = flag.String("h", "185.20.227.83", "Host")
		port     = flag.Int("p", 22, "Port")
	)

	config := &ssh.ClientConfig{
		User:            *user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
	}

	client, err := ssh.Dial("tcp", *host+":"+strconv.Itoa(*port), config)
	if err != nil {
		log.Fatal("Failed to dial: %s", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: %s", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
		if _, err = fmt.Fprintf(stdin, "%s\n", cmd); err != nil {
			log.Fatal(err)
		}
		if bytes.Compare([]byte(cmd), []byte("exit")) == 0 {
			break
		}

	}

	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
