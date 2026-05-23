package email

import (
	"fmt"
	"net/smtp"
	"os"

	"go.uber.org/zap"
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

	zap.L().Info("SMTP sender initialized",
		zap.String("host", sender.host),
		zap.String("port", sender.port),
		zap.String("user", sender.username),
		zap.String("from", sender.from),
	)

	return sender
}

func (s *SMTPSender) SendResetCode(to, code string) error {
	zap.L().Debug("Attempting to send reset code", zap.String("to", to))

	if s.host == "" || s.port == "" || s.username == "" || s.password == "" || s.from == "" {
		err := fmt.Errorf("SMTP configuration incomplete")
		zap.L().Error("SMTP configuration error",
			zap.String("host", s.host),
			zap.String("port", s.port),
			zap.String("user", s.username),
			zap.String("from", s.from),
		)
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

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg)); err != nil {
		zap.L().Error("Failed to send email",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	zap.L().Info("Reset code sent successfully", zap.String("to", to))
	return nil
}
