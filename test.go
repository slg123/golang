package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"gopkg.in/gomail.v2"
)

func getDiskSpace(path string) (total, free int, err error) {
	s := syscall.Statfs_t{}
	err = syscall.Statfs(path, &s)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	total = int(s.Bsize) * int(s.Blocks)
	free  = int(s.Bsize) * int(s.Bfree)
	return
}

func getNoreplyPassword() string {
	file, err := ioutil.ReadFile("/home/zadmin/bin/.noreplypw")
	if err != nil {
		fmt.Println(err.Error())
	}
	pw := string(file)
	return pw
}

func rootIsFull(THRESHHOLD float64) bool {
	total, free, _ := getDiskSpace("/")
	pctUsed := float64(100.00) - ( float64(free) / float64(total) * 100.0 )
	if pctUsed >= THRESHHOLD {
		fmt.Printf("percent used for root filesystem, /, on sbjenkins: %.2f\n", pctUsed)
		return true
	} 
	return false
}

func jenkinsIsDown() bool {
	cmd := "ps -ef|grep [j]enkinsci >/dev/null 2>&1; echo $?"
	status, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		fmt.Sprintf("Failed to execute command: %s", cmd)
	}

	i, _ := strconv.Atoi(strings.Trim(string(status), "\n"))
	// grep failed if $? is not equal to 0 from cmd above then process is DOWN
	if i != 0 { 
		fmt.Printf("jenkinsci process DOWN on sbjenkins! Restart and investigate.\n")
		return true
	} 
	return false
}

type email struct {
	From              string
	NoReplyAcct       string
	To                string
	ToAcct1           string
	ToAcct2           string
	ToAcct3           string
	Subject           string
	TxtHTMLBody       string
	RootFullMsg       string
	RootFullBody2     string
	JenkinsDownMsg    string
	JenkinsDownBody1  string
	JenkinsDownBody2  string
	SMTPServer        string
	SMTPPort          int
}

func main() {

	e := email{}

	e.From		= "From"
	e.NoReplyAcct   = "noreply@r1soft.com"
	e.To		= "To"
	e.ToAcct1       = "scott.gillespie@r1soft.com"
	e.ToAcct2       = "alex.vongluck@r1soft.com"
	e.ToAcct3       = "stan.love@r1soft.com"
	e.TxtHTMLBody   = "text/html"
	e.Subject       = "Subject"
	e.SMTPServer    = "smtp.office365.com"
	e.SMTPPort      = 587

	m := gomail.NewMessage()

	m.SetHeader(e.From, e.NoReplyAcct)
	//m.SetHeader(e.To, e.ToAcct1, e.ToAcct2, e.ToAcct3)
	m.SetHeader(e.To, e.ToAcct1)

	pw := strings.Trim(getNoreplyPassword(), "\n")

	THRESHHOLD := float64(85.0)

	if (rootIsFull(THRESHHOLD)) {
		fmt.Printf("Jenkins root directory is FULL.\n")
		fmt.Printf("sending email...\n")

		e.RootFullMsg    = "SBJENKINS root filesystem full."
		e.RootFullBody2  = "jenkins-root filesystem is full. Clean old build areas."

		m.SetHeader(e.Subject, e.RootFullMsg)
		m.SetBody(e.TxtHTMLBody, e.RootFullBody2)

		d := gomail.NewDialer(e.SMTPServer, e.SMTPPort, e.NoReplyAcct, pw)

		err := d.DialAndSend(m)
		if err != nil {
			log.Fatal(err)
		}
	}

	if (jenkinsIsDown()) {
		fmt.Printf("Jenkins is DOWN.\n")
		fmt.Printf("sending email...\n")

		e.JenkinsDownMsg    = "SBJENKINS JENKINS is DOWN."
		e.JenkinsDownBody2  = "Jenkins is DOWN. Please restart & investigate."

		m.SetHeader(e.Subject, e.JenkinsDownMsg)
		m.SetBody(e.TxtHTMLBody, e.JenkinsDownBody2)

		d := gomail.NewDialer(e.SMTPServer, e.SMTPPort, e.NoReplyAcct, pw)

		err := d.DialAndSend(m)
		if err != nil {
			log.Fatal(err)
		}
	}
}
