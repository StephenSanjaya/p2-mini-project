package helpers

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(email, subject, content string) {
	senderName := os.Getenv("CONFIG_SENDER_NAME")
	port, _ := strconv.Atoi((os.Getenv("CONFIG_SMTP_PORT")))
	host := os.Getenv("CONFIG_SMTP_HOST")
	username := os.Getenv("CONFIG_AUTH_EMAIL")
	password := os.Getenv("CONFIG_AUTH_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", senderName)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendSuccessRegister(email string, url string) {
	SendMail(
		email,
		"Register success",
		fmt.Sprintf("Register success with url: <b>%s<b>", url),
	)
}

func SendSuccessRental(email string, url string) {
	SendMail(
		email,
		"Rental success",
		fmt.Sprintf("invoice rental url: <b>%s<b>", url),
	)
}

func SendSuccessTopUp(email string, url string) {
	SendMail(
		email,
		"Topup success",
		fmt.Sprintf("invoice top up url: <b>%s<b>", url),
	)
}

func SendSuccessPayment(email string, total_price float64) {
	SendMail(
		email,
		"Payment success",
		fmt.Sprintf("success paid with amount of <b>Rp. %.2f<b>", total_price),
	)
}
