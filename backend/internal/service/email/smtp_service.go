package email

type SMTPService struct {
	// настройки smtp
}

func NewSMTPService() *SMTPService {
	return &SMTPService{}
}

func (s *SMTPService) SendConfirmationEmail(to, token, link string) error {
	// TODO: Реальная отправка через SMTP
	// Пока просто логируем
	println("SENDING EMAIL TO:", to)
	println("TOKEN:", token)
	println("LINK:", link)
	return nil
}
