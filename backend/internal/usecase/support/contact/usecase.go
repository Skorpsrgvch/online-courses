package contact

import (
	"context"
	"errors"
	"net/mail"
	"strings"
)

type Input struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Question string `json:"question"`
}

type Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type EmailSender interface {
	SendSupportRequest(name, replyTo, question string) error
}

type Usecase struct {
	emailSender EmailSender
}

func NewUsecase(emailSender EmailSender) (*Usecase, error) {
	if emailSender == nil {
		return nil, errors.New("email sender is required")
	}

	return &Usecase{emailSender: emailSender}, nil
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)
	question := strings.TrimSpace(input.Question)

	if name == "" {
		return nil, errors.New("укажите имя")
	}
	if email == "" {
		return nil, errors.New("укажите email")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("укажите корректный email")
	}
	if question == "" {
		return nil, errors.New("напишите вопрос")
	}
	if len(question) > 5000 {
		return nil, errors.New("вопрос слишком длинный")
	}

	if err := u.emailSender.SendSupportRequest(name, email, question); err != nil {
		return nil, err
	}

	_ = ctx
	return &Output{
		Success: true,
		Message: "Обращение отправлено",
	}, nil
}
