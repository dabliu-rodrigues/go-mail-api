package mail

import (
	"crypto/tls"
	"emailn/internal/domain/campaign"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(campaign *campaign.Campaign) error {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	var emails []string
	for _, contact := range campaign.Contacts {
		emails = append(emails, contact.Email)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", emails...)
	m.SetHeader("Subject", campaign.Name)
	m.SetBody("text/html", campaign.Content)

	return d.DialAndSend(m)
}
