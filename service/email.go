package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService interface {
	SendOTP(to string, otp string) error
}

type emailService struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewEmailService() EmailService {
	return &emailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("EMAIL_FROM"),
	}
}

func (s *emailService) SendOTP(to string, otp string) error {
	subject := "Subject: Registrasi OTP - Go Tiket Konser\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf("<h3>Halo!</h3><p>Terima kasih telah mendaftar di Go Tiket Konser.</p><p>Kode OTP Anda adalah: <strong>%s</strong></p><p>Kode ini berlaku selama 5 menit.</p>", otp)
	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
