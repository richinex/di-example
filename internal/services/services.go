package services

import (
    "fmt"
    "di-example/pkg/logger"
)

// Interfaces remain the same
type UserService interface {
    GetUser(id int) string
}

type EmailService interface {
    SendEmail(to, message string) error
}

type ConfigService interface {
    GetConfig() string
}

// UserService implementation
type userService struct {
    prefix string
}

func NewUserService() UserService {
    log := logger.Get()
    log.Infow("Creating new UserService", "prefix", "USER-")
    return &userService{prefix: "USER-"}
}

func (s *userService) GetUser(id int) string {
    log := logger.Get()
    result := fmt.Sprintf("%s%d", s.prefix, id)
    log.Infow("Getting user",
        "id", id,
        "prefix", s.prefix,
        "result", result)
    return result
}

// EmailService implementation
type emailService struct {
    server string
}

func NewEmailService() EmailService {
    log := logger.Get()
    log.Infow("Creating new EmailService", "server", "smtp.example.com")
    return &emailService{server: "smtp.example.com"}
}

func (s *emailService) SendEmail(to, message string) error {
    log := logger.Get()
    log.Infow("Sending email",
        "to", to,
        "server", s.server,
        "messageLength", len(message))

    fmt.Printf("Sending email to %s via %s: %s\n", to, s.server, message)

    log.Infow("Email sent successfully",
        "to", to,
        "server", s.server)
    return nil
}

// ConfigService implementation
type configService struct {
    env string
}

func NewConfigService() ConfigService {
    log := logger.Get()
    log.Infow("Creating new ConfigService", "environment", "development")
    return &configService{env: "development"}
}

func (s *configService) GetConfig() string {
    log := logger.Get()
    result := fmt.Sprintf("Environment: %s", s.env)
    log.Infow("Getting config",
        "environment", s.env,
        "result", result)
    return result
}