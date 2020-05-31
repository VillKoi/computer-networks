package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
)

const (
	MAIL_TEMPLATE = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-Version: 1.0
Content-Type: text/plain; charset="utf-8"
{{.Body}}
`
)

type conf struct {
	smtphost, serveport, user, password string
}

type Message struct {
	From, To, Subject, tplname string
}

type SmtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

var (
	cnf conf
)

func init() {
	file, err := ioutil.ReadFile("../params.txt")
	if err != nil {
		log.Fatal(err)
	}

	var c map[string]interface{}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal(err)
	}

	password, _ := base64.StdEncoding.DecodeString(c["password"].(string))

	cnf = conf{
		c["server"].(string),
		c["port"].(string),
		c["user"].(string),
		string(password),
	}
}

func parseTemplate(from, to, subject, body string) ([]byte, error) {

	var (
		err error
		doc bytes.Buffer
	)

	context := &SmtpTemplateData{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}
	parsed := template.New("emailTemplate")

	parsed, err = parsed.Parse(MAIL_TEMPLATE)
	if err != nil {
		log.Print("error trying to parse mail template")
	}
	err = parsed.Execute(&doc, context)
	return doc.Bytes(), err
}

func getSMTPClient() *smtp.Client {
	client, err := smtp.Dial(cnf.smtphost + ":" + cnf.serveport)
	if err != nil {
		log.Fatal("dial", err)
	}

	tlc := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         cnf.smtphost,
	}

	if err := client.StartTLS(tlc); err != nil {
		log.Println("tls error: ", err)
	}

	auth := smtp.PlainAuth("", cnf.user, cnf.password, cnf.smtphost)
	if err = client.Auth(auth); err != nil {
		log.Println("new client", err)
	}

	return client
}

func main() {
	client := getSMTPClient()
	defer client.Quit()
	var from, to, subject string

	for {
		fmt.Println("Enter From: ")
		fmt.Scan(&from)

		fmt.Println("Enter To: ")
		fmt.Scan(&to)

		fmt.Println("Enter Subject: ")
		fmt.Scan(&subject)

		fmt.Println(`Enter message and "end" at the end: `)

		var body string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			str := scanner.Text()
			if err := scanner.Err(); err != nil {
				log.Println("Scanner: ", err)
			}

			if bytes.Compare([]byte(str), []byte("end")) == 0 {
				fmt.Println()
				break
			}
			body += str + "\n"
		}

		log.Println(string(body))
		content, err := parseTemplate(from, to, subject, body)
		err = client.Noop()
		if err != nil {
			log.Println("reestablish connection", err)
			client = getSMTPClient()
		}

		if err = client.Mail(from); err != nil {
			log.Println(err)
		}

		if err = client.Rcpt(to); err != nil {
			log.Println(err)
		}

		writecloser, err := client.Data()
		if err != nil {
			log.Println(err)
		}

		_, err = writecloser.Write(content)
		if err != nil {
			log.Println(err)
		}

		if err := writecloser.Close(); err != nil {
			log.Println("wc close: ", err)
		}

		fmt.Println("Message sent\nContinue? (yes/no)\n")
		var boo string
		fmt.Scan(&boo)
		if boo == "no" {
			break
		}

	}
}
