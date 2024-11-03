package container

import (
    "fmt"
    "reflect"
    "sync"
)

// Container represents a dependency injection container
type Container struct {
    mu       sync.RWMutex
    services map[string]interface{}
}

// NewContainer creates a new DI container instance
func NewContainer() *Container {
    return &Container{
        services: make(map[string]interface{}),
    }
}

// Register adds a service to the container with the given qualifier
func (c *Container) Register(qualifier string, service interface{}) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if service == nil {
        return fmt.Errorf("cannot register nil service for qualifier: %s", qualifier)
    }

    // Check if service is already registered
    if _, exists := c.services[qualifier]; exists {
        return fmt.Errorf("service already registered for qualifier: %s", qualifier)
    }

    c.services[qualifier] = service
    return nil
}

// Resolve retrieves a service from the container by qualifier
func (c *Container) Resolve(qualifier string) (interface{}, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    service, exists := c.services[qualifier]
    if !exists {
        return nil, fmt.Errorf("no service found for qualifier: %s", qualifier)
    }

    return service, nil
}

// InjectStruct injects dependencies into struct fields marked with "di" tags
func (c *Container) InjectStruct(target interface{}) error {
    // Target must be a pointer to struct
    targetValue := reflect.ValueOf(target)
    if targetValue.Kind() != reflect.Ptr {
        return fmt.Errorf("target must be a pointer to struct, got: %v", targetValue.Kind())
    }

    // Get the actual struct value
    targetValue = targetValue.Elem()
    if targetValue.Kind() != reflect.Struct {
        return fmt.Errorf("target must be a pointer to struct, got pointer to: %v", targetValue.Kind())
    }

    targetType := targetValue.Type()

    // Iterate over all fields in the struct
    for i := 0; i < targetType.NumField(); i++ {
        field := targetType.Field(i)

        // Check if field has DI tag
        qualifier, ok := field.Tag.Lookup("di")
        if !ok {
            continue // Skip fields without DI tag
        }

        // Resolve the dependency
        service, err := c.Resolve(qualifier)
        if err != nil {
            return fmt.Errorf("failed to inject field %s: %v", field.Name, err)
        }

        // Get field value
        fieldValue := targetValue.Field(i)
        if !fieldValue.CanSet() {
            return fmt.Errorf("cannot set field %s (unexported)", field.Name)
        }

        // Check if types are compatible
        serviceValue := reflect.ValueOf(service)
        if !serviceValue.Type().AssignableTo(fieldValue.Type()) {
            return fmt.Errorf("service type %v is not assignable to field type %v",
                serviceValue.Type(), fieldValue.Type())
        }

        // Set the field value
        fieldValue.Set(serviceValue)
    }

    return nil
}