package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"syscall"
	"gopkg.in/gomail.v2"
)

func getNoreplyPassword() (pw string) {

	file, err := ioutil.ReadFile(".noreplypw")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(file))
	return pw
}


func main() {

	pw := getNoreplyPassword()

	fmt.Println(pw)

	stat := syscall.Statfs_t{}
	err  := syscall.Statfs("/", &stat)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	n := stat.Bfree - stat.Bavail
	m := stat.Blocks
	p := float64(n) / float64(m) * float64(100.0)

	fmt.Printf("%.2f\n", p)

	threshhold := float64(95.0)

	if p >= threshhold {
		fmt.Printf("sending mail...\n"); 
		m := gomail.NewMessage()
		m.SetHeader("From", "noreply@r1soft.com")
		//m.SetHeader("To", "scott.gillespie@r1soft.com", "alex.vongluck@r1soft.com", "stan.love@r1soft.com")
		m.SetHeader("To", "scott.gillespie@r1soft.com")
		m.SetHeader("Subject", "SBJENKINS root filesystem at or over 95% full.")
		m.SetBody("text/html", "jenkins-root filesystem above 95% used. Please cleanup some old builds, if possible.")

		//d := gomail.NewDialer("smtp.office365.com", 587, "noreply@r1soft.com", "")
		d := gomail.NewDialer("smtp.office365.com", 587, "noreply@r1soft.com", pw)

		err := d.DialAndSend(m)
		if err != nil {
			log.Fatal(err)
		}
	} 
}
