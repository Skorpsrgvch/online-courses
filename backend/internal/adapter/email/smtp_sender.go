package email

import (
	"fmt"
	"mime"
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

	if err := s.validateConfig(); err != nil {
		return err
	}

	subject := "Код восстановления пароля"
	body := fmt.Sprintf(`Здравствуйте!

Ваш код для восстановления пароля: %s

Код действителен в течение 15 минут.

Если вы не запрашивали восстановление, просто проигнорируйте это письмо.

С уважением, команда WomanFormula`, code)

	if err := s.sendPlainText(to, subject, "", body); err != nil {
		zap.L().Error("Failed to send email",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	zap.L().Info("Reset code sent successfully", zap.String("to", to))
	return nil
}

func (s *SMTPSender) SendCourseAccessGranted(to, userName, courseTitle string) error {
	zap.L().Debug("Attempting to send course access notification", zap.String("to", to))

	if err := s.validateConfig(); err != nil {
		return err
	}

	if userName == "" {
		userName = "student"
	}

	subject := "Вам выдан доступ к курсу"
	fromName := s.from
	body := fmt.Sprintf(`Здравствуйте, %s!

Вам выдан доступ к курсу "%s".

Вы можете войти в личный кабинет и начать обучение.

С уважением, команда WomanFormula`, userName, courseTitle)

	if err := s.sendPlainText(to, subject, fromName, body); err != nil {
		zap.L().Error("Failed to send course access notification",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send course access notification: %w", err)
	}

	zap.L().Info("Course access notification sent successfully", zap.String("to", to))
	return nil
}

func (s *SMTPSender) SendSupportRequest(name, replyTo, question string) error {
	zap.L().Debug("Attempting to send support request", zap.String("reply_to", replyTo))

	if err := s.validateConfig(); err != nil {
		return err
	}

	supportEmail := os.Getenv("SUPPORT_EMAIL")
	if supportEmail == "" {
		supportEmail = s.from
	}

	subject := "Новое обращение в поддержку"
	fromName := s.from
	body := fmt.Sprintf(`Новое обращение в поддержку.

Имя: %s
Email: %s

Вопрос:
%s`, name, replyTo, question)

	if err := s.sendPlainText(supportEmail, subject, fromName, body, replyTo); err != nil {
		zap.L().Error("Failed to send support request",
			zap.String("support_email", supportEmail),
			zap.String("reply_to", replyTo),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send support request: %w", err)
	}

	zap.L().Info("Support request sent successfully",
		zap.String("support_email", supportEmail),
		zap.String("reply_to", replyTo),
	)
	return nil
}

func (s *SMTPSender) validateConfig() error {
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

	return nil
}

func (s *SMTPSender) sendPlainText(to, subject, fromName, body string, replyTo ...string) error {
	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)
	fromHeader := s.from
	if fromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("UTF-8", fromName), s.from)
	}

	replyToHeader := ""
	if len(replyTo) > 0 && replyTo[0] != "" {
		replyToHeader = fmt.Sprintf("Reply-To: %s\r\n", replyTo[0])
	}

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"%s"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", fromHeader, to, replyToHeader, encodedSubject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
