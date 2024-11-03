package services

import (
    "fmt"
)

// UserService handles user-related operations
type UserService interface {
    GetUser(id int) string
}

// EmailService handles email operations
type EmailService interface {
    SendEmail(to, message string) error
}

// ConfigService handles configuration
type ConfigService interface {
    GetConfig() string
}

// Implementation of UserService
type userService struct {
    prefix string
}

func NewUserService() UserService {
    return &userService{prefix: "USER-"}
}

func (s *userService) GetUser(id int) string {
    return fmt.Sprintf("%s%d", s.prefix, id)
}

// Implementation of EmailService
type emailService struct {
    server string
}

func NewEmailService() EmailService {
    return &emailService{server: "smtp.example.com"}
}

func (s *emailService) SendEmail(to, message string) error {
    fmt.Printf("Sending email to %s via %s: %s\n", to, s.server, message)
    return nil
}

// Implementation of ConfigService
type configService struct {
    env string
}

func NewConfigService() ConfigService {
    return &configService{env: "development"}
}

func (s *configService) GetConfig() string {
    return fmt.Sprintf("Environment: %s", s.env)
}