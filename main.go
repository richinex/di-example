package main

import (
	"di-example/internal/models"
	"di-example/internal/services"
	"di-example/pkg/container"
	"di-example/pkg/logger"
	"di-example/pkg/reflection"
	"fmt"
)

func main() {
    // Initialize logger
    logger.Initialize(true) // true for development mode with colors
    defer logger.Sync()

    log := logger.Get()
    log.Info("Starting application")

    // Create new DI container
    log.Info("Initializing DI container")
    di := container.NewContainer()

    // Create services
    log.Info("Creating services")
    userService := services.NewUserService()
    emailService := services.NewEmailService()
    configService := services.NewConfigService()

	// Inversion of Control (IoC)
	// The Container manages service lifecycle
	// Services are registered and resolved through the container

    // Register services
    log.Info("Registering services in container")
    if err := di.Register("userService", userService); err != nil {
        log.Fatalw("Failed to register userService", "error", err)
    }
    if err := di.Register("emailService", emailService); err != nil {
        log.Fatalw("Failed to register emailService", "error", err)
    }
    if err := di.Register("configService", configService); err != nil {
        log.Fatalw("Failed to register configService", "error", err)
    }

    // Create injectable struct
    log.Info("Creating injectable struct")
    injectable := &models.Injectable{}

    // Inject dependencies
	// Dependency Injection
	// Field injection through struct tags
	// Resolution of dependencies
    log.Info("Injecting dependencies")
    if err := di.InjectStruct(injectable); err != nil {
        log.Fatalw("Failed to inject dependencies", "error", err)
    }

    // Create reflection inspector
    log.Info("Creating reflection inspector")
    inspector := reflection.NewInspector()

    // Inspect the injectable struct
    log.Info("Inspecting injectable struct")
    info, err := inspector.InspectStruct(injectable)
    if err != nil {
        log.Fatalw("Failed to inspect struct", "error", err)
    }

    // Print inspection results
    log.Info("=== Struct Inspection Results ===")
    fmt.Println(inspector.PrettyPrint(info))

    // Test services
    log.Info("=== Testing Injected Services ===")

    if us, ok := injectable.UserService.(services.UserService); ok {
        result := us.GetUser(123)
        log.Infow("Tested UserService", "result", result)
    }

    if es, ok := injectable.EmailService.(services.EmailService); ok {
        err := es.SendEmail("test@example.com", "Hello from DI!")
        log.Infow("Tested EmailService", "error", err)
    }

    if cs, ok := injectable.ConfigService.(services.ConfigService); ok {
        result := cs.GetConfig()
        log.Infow("Tested ConfigService", "result", result)
    }

    log.Info("Application completed successfully")
}