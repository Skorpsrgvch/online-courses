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
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f9fafb;font-family:'Helvetica Neue',Arial,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f9fafb;padding:40px 0;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.06);">
    
    <tr><td style="background:linear-gradient(135deg,#f43f5e,#e11d48);padding:32px 32px 28px;text-align:center;">
        <p style="margin:0;font-size:26px;font-weight:300;color:#ffffff;letter-spacing:0.5px;">Тазовое здоровье</p>
    </td></tr>
    
    <tr><td style="padding:36px 32px 20px;">
        <p style="margin:0 0 8px;font-size:18px;font-weight:600;color:#1f2937;">Здравствуйте!</p>
        <p style="margin:0 0 24px;font-size:15px;line-height:1.6;color:#4b5563;">
            Вы запросили восстановление пароля. Скопируйте код ниже и вставьте его в форму на сайте:
        </p>
        
        <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="background-color:#fff1f2;border:2px dashed #f43f5e;border-radius:12px;padding:22px 16px;text-align:center;">
            <span style="font-family:'Courier New',monospace;font-size:32px;font-weight:700;letter-spacing:8px;color:#e11d48;">%s</span>
        </td></tr>
        </table>
        
        <p style="margin:24px 0 0;font-size:13px;line-height:1.5;color:#9ca3af;">
            Код действителен в течение <strong style="color:#6b7280;">15 минут</strong>.<br>
            Если вы не запрашивали восстановление пароля, просто проигнорируйте это письмо.
        </p>
    </td></tr>
    
    <tr><td style="background-color:#f9fafb;padding:20px 32px;text-align:center;border-top:1px solid #f3f4f6;">
        <p style="margin:0;font-size:12px;color:#9ca3af;">
            С уважением, <strong style="color:#e11d48;">Ольга Пимченко</strong>
        </p>
    </td></tr>
    
</table>
</td></tr>
</table>
</body>
</html>`, code)

	if err := s.sendHTML(to, subject, "Тазовое здоровье", body); err != nil {
		zap.L().Error("Failed to send email", zap.String("to", to), zap.Error(err))
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
		userName = "клиент"
	}

	subject := "Вам открыт доступ к курсу"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f9fafb;font-family:'Helvetica Neue',Arial,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f9fafb;padding:40px 0;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.06);">
    
    <tr><td style="background:linear-gradient(135deg,#f43f5e,#e11d48);padding:32px 32px 28px;text-align:center;">
        <p style="margin:0;font-size:26px;font-weight:300;color:#ffffff;letter-spacing:0.5px;">Тазовое здоровье</p>
    </td></tr>
    
    <tr><td style="padding:36px 32px 20px;">
        <p style="margin:0 0 8px;font-size:18px;font-weight:600;color:#1f2937;">Здравствуйте, %s!</p>
        <p style="margin:0 0 12px;font-size:15px;line-height:1.6;color:#4b5563;">
            Хорошие новости — вам открыт доступ к курсу:
        </p>
        
        <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="background-color:#fff1f2;border-radius:12px;padding:20px 16px;text-align:center;">
            <span style="font-size:20px;font-weight:600;color:#e11d48;">«%s»</span>
        </td></tr>
        </table>
        
        <p style="margin:24px 0 0;font-size:15px;line-height:1.6;color:#4b5563;">
            Вы можете войти в личный кабинет и приступить к обучению прямо сейчас.
        </p>
        
        <table width="100%%" cellpadding="0" cellspacing="0" style="margin-top:24px;">
        <tr><td align="center">
            <a href="https://olga-pimchenko.ru/login" target="_blank" style="display:inline-block;background:linear-gradient(135deg,#f43f5e,#e11d48);color:#ffffff;text-decoration:none;padding:14px 40px;border-radius:10px;font-size:15px;font-weight:600;">Перейти к обучению</a>
        </td></tr>
        </table>
    </td></tr>
    
    <tr><td style="background-color:#f9fafb;padding:20px 32px;text-align:center;border-top:1px solid #f3f4f6;">
        <p style="margin:0;font-size:12px;color:#9ca3af;">
            С уважением,  <strong style="color:#e11d48;">Ольга Пимченко</strong>
        </p>
    </td></tr>
    
</table>
</td></tr>
</table>
</body>
</html>`, userName, courseTitle)

	if err := s.sendHTML(to, subject, "Тазовое здоровье", body); err != nil {
		zap.L().Error("Failed to send course access notification", zap.String("to", to), zap.Error(err))
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
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f9fafb;font-family:'Helvetica Neue',Arial,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f9fafb;padding:40px 0;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.06);">
    
    <tr><td style="background:linear-gradient(135deg,#f43f5e,#e11d48);padding:32px 32px 28px;text-align:center;">
        <p style="margin:0;font-size:26px;font-weight:300;color:#ffffff;letter-spacing:0.5px;">Новое обращение</p>
    </td></tr>
    
    <tr><td style="padding:36px 32px 20px;">
        <table width="100%%" cellpadding="0" cellspacing="0" style="margin-bottom:20px;">
            <tr><td style="padding:8px 0;font-size:13px;color:#9ca3af;width:60px;">Имя:</td>
                <td style="padding:8px 0;font-size:15px;color:#1f2937;font-weight:500;">%s</td></tr>
            <tr><td style="padding:8px 0;font-size:13px;color:#9ca3af;">Email:</td>
                <td style="padding:8px 0;font-size:15px;color:#e11d48;">%s</td></tr>
        </table>
        
        <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="background-color:#f9fafb;border-radius:12px;padding:20px 16px;">
            <p style="margin:0 0 8px;font-size:13px;color:#9ca3af;text-transform:uppercase;letter-spacing:1px;">Вопрос</p>
            <p style="margin:0;font-size:15px;line-height:1.6;color:#1f2937;">%s</p>
        </td></tr>
        </table>
    </td></tr>
    
</table>
</td></tr>
</table>
</body>
</html>`, name, replyTo, question)

	if err := s.sendHTML(supportEmail, subject, "Поддержка", body, replyTo); err != nil {
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

// validateConfig проверяет, что все поля конфигурации заполнены
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

// sendHTML отправляет HTML-письмо
func (s *SMTPSender) sendHTML(to, subject, fromName, htmlBody string, replyTo ...string) error {
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
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", fromHeader, to, replyToHeader, encodedSubject, htmlBody)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
