package helpers

// func SendMail(email, subject, content string) {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", "example@hacktiv8.com")
// 	m.SetHeader("To", email)
// 	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/html", content)

// 	d := gomail.NewDialer(
// 		"smtp-example.com",
// 		587,
// 		"example",
// 		"example",
// 	)

// 	// Send the email to Bob, Cora and Dan.
// 	if err := d.DialAndSend(m); err != nil {
// 		panic(err)
// 	}
// }

// func SendSuccessCreateRent(email string) {
// 	SendMail(
// 		email,
// 		"example subject",
// 		"example content",
// 	)
// }
