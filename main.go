package main

import (
    "fmt"
    "log"

    "di-example/internal/models"
    "di-example/internal/services"
    "di-example/pkg/container"
    "di-example/pkg/reflection"
)

func main() {
    // Create new DI container
    di := container.NewContainer()

    // Create services
    userService := services.NewUserService()
    emailService := services.NewEmailService()
    configService := services.NewConfigService()

    // Register services in container
    if err := di.Register("userService", userService); err != nil {
        log.Fatalf("Failed to register userService: %v", err)
    }
    if err := di.Register("emailService", emailService); err != nil {
        log.Fatalf("Failed to register emailService: %v", err)
    }
    if err := di.Register("configService", configService); err != nil {
        log.Fatalf("Failed to register configService: %v", err)
    }

    // Create injectable struct
    injectable := &models.Injectable{}

    // Inject dependencies
    if err := di.InjectStruct(injectable); err != nil {
        log.Fatalf("Failed to inject dependencies: %v", err)
    }

    // Create reflection inspector
    inspector := reflection.NewInspector()

    // Inspect the injectable struct
    info, err := inspector.InspectStruct(injectable)
    if err != nil {
        log.Fatalf("Failed to inspect struct: %v", err)
    }

    // Print inspection results
    fmt.Println("=== Struct Inspection Results ===")
    fmt.Println(inspector.PrettyPrint(info))

    // Demonstrate usage of injected services
    fmt.Println("\n=== Testing Injected Services ===")

    // Test UserService
    if us, ok := injectable.UserService.(services.UserService); ok {
        fmt.Printf("UserService result: %s\n", us.GetUser(123))
    }

    // Test EmailService
    if es, ok := injectable.EmailService.(services.EmailService); ok {
        es.SendEmail("test@example.com", "Hello from DI!")
    }

    // Test ConfigService
    if cs, ok := injectable.ConfigService.(services.ConfigService); ok {
        fmt.Printf("ConfigService result: %s\n", cs.GetConfig())
    }
}