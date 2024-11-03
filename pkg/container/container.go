// Package container provides dependency injection functionality
package container

import (
    "fmt"
    "reflect"
    "sync"
    "di-example/pkg/logger"
    "go.uber.org/zap"
)

// Container represents a dependency injection container that manages services
type Container struct {
    mu       sync.RWMutex                // Mutex for thread-safe operations
    services map[string]interface{}      // Map to store services with their qualifiers
    log      *zap.SugaredLogger         // Logger instance
}

// NewContainer creates and initializes a new DI container
func NewContainer() *Container {
    return &Container{
        services: make(map[string]interface{}), // Initialize empty service map
        log:      logger.Get(),                 // Get logger instance
    }
}

// Register adds a new service to the container with the specified qualifier
func (c *Container) Register(qualifier string, service interface{}) error {
    c.mu.Lock()                    // Lock for thread safety
    defer c.mu.Unlock()            // Ensure unlock when function returns

    // Log registration attempt
    c.log.Infow("Registering service",
        "qualifier", qualifier,
        "type", reflect.TypeOf(service))

    // Validate service is not nil
    if service == nil {
        c.log.Errorw("Cannot register nil service",
            "qualifier", qualifier)
        return fmt.Errorf("cannot register nil service for qualifier: %s", qualifier)
    }

    // Check if service already exists
    if _, exists := c.services[qualifier]; exists {
        c.log.Errorw("Service already registered",
            "qualifier", qualifier)
        return fmt.Errorf("service already registered for qualifier: %s", qualifier)
    }

    // Store service in container
    c.services[qualifier] = service
    c.log.Infow("Service registered successfully",
        "qualifier", qualifier,
        "type", reflect.TypeOf(service))
    return nil
}

// Resolve retrieves a service from the container by its qualifier
func (c *Container) Resolve(qualifier string) (interface{}, error) {
    c.mu.RLock()                   // Read lock for thread safety
    defer c.mu.RUnlock()           // Ensure unlock when function returns

    c.log.Debugw("Resolving service", "qualifier", qualifier)

    // Look up service in container
    service, exists := c.services[qualifier]
    if !exists {
        c.log.Errorw("Service not found", "qualifier", qualifier)
        return nil, fmt.Errorf("no service found for qualifier: %s", qualifier)
    }

    c.log.Debugw("Service resolved successfully",
        "qualifier", qualifier,
        "type", reflect.TypeOf(service))
    return service, nil
}

// InjectStruct injects dependencies into struct fields marked with "di" tags
func (c *Container) InjectStruct(target interface{}) error {
    c.log.Info("Starting struct injection")

    // Get reflect.Value of target and ensure it's a pointer
    targetValue := reflect.ValueOf(target)
    if targetValue.Kind() != reflect.Ptr {
        c.log.Errorw("Target must be a pointer",
            "actualKind", targetValue.Kind())
        return fmt.Errorf("target must be a pointer to struct, got: %v", targetValue.Kind())
    }

    // Dereference pointer to get struct value
    targetValue = targetValue.Elem()
    targetType := targetValue.Type()

    // Verify target is a struct
    if targetValue.Kind() != reflect.Struct {
        c.log.Errorw("Target must be a pointer to struct",
            "actualKind", targetValue.Kind())
        return fmt.Errorf("target must be a pointer to struct, got pointer to: %v", targetValue.Kind())
    }

    c.log.Infow("Analyzing struct for injection",
        "structType", targetType.Name(),
        "numFields", targetType.NumField())

    // Iterate through all fields in the struct
    for i := 0; i < targetType.NumField(); i++ {
        field := targetType.Field(i)

        // Look for 'di' tag on field
        qualifier, ok := field.Tag.Lookup("di")
        if !ok {
            c.log.Debugw("Skipping field without di tag",
                "field", field.Name)
            continue
        }

        c.log.Infow("Injecting field",
            "field", field.Name,
            "qualifier", qualifier)

        // Get field value and check if it can be set
        fieldValue := targetValue.Field(i)
        if !fieldValue.CanSet() {
            c.log.Debugw("Cannot set field (unexported), skipping",
                "field", field.Name)
            continue
        }

        // Resolve service for this field
        service, err := c.Resolve(qualifier)
        if err != nil {
            // If the service is not found, just log it and continue
            c.log.Debugw("Optional service not found, skipping field",
                "field", field.Name,
                "qualifier", qualifier)
            continue
        }

        // Verify type compatibility
        serviceValue := reflect.ValueOf(service)
        if !serviceValue.Type().AssignableTo(fieldValue.Type()) {
            c.log.Errorw("Type mismatch during injection",
                "field", field.Name,
                "expectedType", fieldValue.Type(),
                "actualType", serviceValue.Type())
            return fmt.Errorf("service type %v is not assignable to field type %v",
                serviceValue.Type(), fieldValue.Type())
        }

        // Set the field value to the service
        fieldValue.Set(serviceValue)
        c.log.Infow("Successfully injected field",
            "field", field.Name,
            "qualifier", qualifier)
    }

    c.log.Info("Completed struct injection")
    return nil
}