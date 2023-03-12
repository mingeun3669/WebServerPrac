package libs

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
)

type Uuid map[string]string

func SendMail(email string) (string, error) {
	auth := smtp.PlainAuth("", "sdh220312@sdh.hs.kr", "smzqtywdkrocvxui", "smtp.gmail.com")
	from := "sdh220312@sdh.hs.kr"
	to := []string{email}
	num := ""
	min := 65
	max := 90
	for i := 0; i < 6; i++ {
		a := rand.Intn(max-min) + min
		num += string(a)
	}

	headerSubject := "Subject: Verify Mail\r\n"
	headerBlank := "\r\n"
	body := "please input this number. " + num
	msg := []byte(headerSubject + headerBlank + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("success")
	return num, nil
}

func (u Uuid) Add(email, num string) error {
	_, err := u.Search(email)
	if err != nil {
		u[email] = num
	} else {
		return errors.New("Already Existed")
	}
	return nil
}

func (u Uuid) Search(email string) (string, error) {
	num, ex := u[email]
	if ex {
		return num, nil
	}
	return "", errors.New("Not Found")
}
