package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gliderlabs/ssh"
)

func parseCommand(s string) []string {
	s = strings.TrimRight(s, "\n")
	re := regexp.MustCompile(`\s+`)
	re.ReplaceAllString(s, " ")
	data := strings.Split(s, " ")
	return data
}

func main()  {

	var (
		user = flag.String("user", "iu9", "Username for login")
		pass = flag.String("pass", "qwerty", "Password for login")
		port = flag.Int("port", 22, "Port")
		)


	flag.Parse()

	server := &ssh.Server{
		Addr: fmt.Sprintf(":%d", *port),
		Handler: func(std ssh.Session) {
			io.WriteString(std, fmt.Sprintf("You've been connected to %s\n", std.LocalAddr().String()))
			reader := bufio.NewReader(std)
		loop:
			for {
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("GetLines: " + err.Error())
					break
				}
				fmt.Println(text)

				command := parseCommand(text)

				switch command[0] {
				case "exit":
					break loop
				case "cd":
					if len(command) < 2 {
						home := os.Getenv("HOME")
						os.Chdir(home)
					} else {
						err := os.Chdir(command[1])
						if err != nil {
							io.WriteString(std, err.Error())
						}
					}
				case "mkdir":
					err := os.Mkdir(command[1], os.ModePerm)
					if err != nil {
						io.WriteString(std, err.Error())
					}
					io.WriteString(std, "\n")
				case "ls":
					files, err := ioutil.ReadDir(".")
					if err != nil {
						io.WriteString(std, err.Error())
					}

					for _, file := range files {
						fmt.Fprint(std, file.Name(), "\n")
					}
					fmt.Fprint(std, "\n")
				case "rm":
					if len(command) < 3 || command[1] != "-r" {
						fmt.Fprintf(std, "Command %s not find, did you mean rm -r namedir", text)
					} else {
						err = os.RemoveAll(command[2])
						if err != nil {
							io.WriteString(std, err.Error())
						}
					}

				case "echo":
					if len(command) < 2 {
						io.WriteString(std, "\n")
					} else {
						io.WriteString(std, strings.Join(command[1:], " ")+"\n")
					}
				default:
					out, err := exec.Command(command[0], command[1:]...).Output()
					if err != nil {
						io.WriteString(std, err.Error())
					}
					io.WriteString(std, fmt.Sprintf("%s\n", string(out)))
				}
			}

			err := std.Exit(0)
			if err != nil {
				fmt.Println(err)
			}
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			return ctx.User() == *user && password == *pass
		},
		PtyCallback: func(ctx ssh.Context, pty ssh.Pty) bool {
			return false
		},
	}

	log.Println(fmt.Sprintf("starting ssh server on port %d...", *port))
	log.Fatal(server.ListenAndServe())

}