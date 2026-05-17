package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPSender() *SMTPSender {
	sender := &SMTPSender{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
		from:     os.Getenv("EMAIL_FROM"),
	}

	// Логирование конфигурации (пароль не пишем!)
	log.Printf("[SMTP] Initialized: Host=%s, Port=%s, User=%s, From=%s",
		sender.host, sender.port, sender.username, sender.from)

	return sender
}

func (s *SMTPSender) SendResetCode(to, code string) error {
	log.Printf("[SMTP] Attempting to send code to: %s", to)

	// Проверка заполненности полей
	if s.host == "" || s.port == "" || s.username == "" || s.password == "" || s.from == "" {
		err := fmt.Errorf("SMTP configuration incomplete: host=%s, port=%s, user=%s, from=%s",
			s.host, s.port, s.username, s.from)
		log.Printf("[SMTP] Error: %v", err)
		return err
	}

	subject := "Код восстановления пароля"
	body := fmt.Sprintf(`Здравствуйте!

Ваш код для восстановления пароля: %s

Код действителен в течение 15 минут.

Если вы не запрашивали восстановление, просто проигнорируйте это письмо.

С уважением, команда WomanFormula`, code)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	log.Printf("[SMTP] Connecting to %s...", addr)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("[SMTP] Failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[SMTP] Successfully sent code to %s", to)
	return nil
}
