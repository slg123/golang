package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"gopkg.in/gomail.v2"
)

func getNoreplyPassword() string {
	file, err := ioutil.ReadFile("/home/continuum/bin/.noreplypw")
	if err != nil {
		fmt.Println(err.Error())
	}
	pw := string(file)
	return pw
}

func proftpdDown() bool {
	cmd := "ps -ef|grep [p]roftpd >/dev/null 2>&1; echo $?"
	status, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		fmt.Sprintf("Failed to execute command: %s", cmd)
	}

	i, _ := strconv.Atoi(strings.Trim(string(status), "\n"))
	if i != 0 { 
		return true
	} 
	return false
}

type email struct {
	From              string
	NoReplyAcct       string
	To                string
	ToAcct            string
	Subject           string
	TxtHTMLBody       string
	proftpdDownMsg    string
	proftpDownBody1   string
	proftpdDownBody2  string
	SMTPServer        string
	SMTPPort          int
}

func main() {

	e := email{}

	e.From		= "From"
	e.NoReplyAcct   = "noreply@r1soft.com"
	e.To		= "To"
	e.ToAcct        = "c247devops@r1soft.com"
	e.TxtHTMLBody   = "text/html"
	e.Subject       = "Subject"
	e.SMTPServer    = "smtp.office365.com"
	e.SMTPPort      = 587

	m := gomail.NewMessage()

	m.SetHeader(e.From, e.NoReplyAcct)
	m.SetHeader(e.To, e.ToAcct)

	pw := strings.Trim(getNoreplyPassword(), "\n")

	if (proftpdDown()) {

		hostname, _ := os.Hostname()

		e.proftpdDownMsg    = "PROFTPD DOWN ON " + hostname
		e.proftpdDownBody2  = "proftpd server is down on host: " + hostname

		m.SetHeader(e.Subject, e.proftpdDownMsg)
		m.SetBody(e.TxtHTMLBody, e.proftpdDownBody2)

		d := gomail.NewDialer(e.SMTPServer, e.SMTPPort, e.NoReplyAcct, pw)

		err := d.DialAndSend(m)
		if err != nil {
			log.Fatal(err)
		}
	}
}
