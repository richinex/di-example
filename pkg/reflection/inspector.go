package reflection

import (
    "fmt"
    "reflect"
    "strings"
)

// StructInfo holds information about a struct
type StructInfo struct {
    Name   string
    Fields []FieldInfo
}

// FieldInfo holds information about a struct field
type FieldInfo struct {
    Name      string
    Type      string
    Tags      map[string]string
    Value     interface{}
    IsExported bool
}

// Inspector provides methods to inspect structs using reflection
type Inspector struct{}

// NewInspector creates a new Inspector instance
func NewInspector() *Inspector {
    return &Inspector{}
}

// InspectStruct analyzes a struct and returns detailed information about it
func (i *Inspector) InspectStruct(target interface{}) (*StructInfo, error) {
    targetValue := reflect.ValueOf(target)
    targetType := targetValue.Type()

    // Handle pointer types
    if targetType.Kind() == reflect.Ptr {
        targetValue = targetValue.Elem()
        targetType = targetValue.Type()
    }

    // Ensure we're dealing with a struct
    if targetType.Kind() != reflect.Struct {
        return nil, fmt.Errorf("target must be a struct or pointer to struct, got: %v",
            targetType.Kind())
    }

    info := &StructInfo{
        Name:   targetType.Name(),
        Fields: make([]FieldInfo, 0, targetType.NumField()),
    }

    // Analyze each field
    for i := 0; i < targetType.NumField(); i++ {
        field := targetType.Field(i)
        fieldValue := targetValue.Field(i)

        // Parse field tags
        tags := make(map[string]string)
        for _, tag := range strings.Split(string(field.Tag), " ") {
            parts := strings.Split(tag, ":")
            if len(parts) == 2 {
                key := strings.Trim(parts[0], "`")
                value := strings.Trim(parts[1], `"`)
                tags[key] = value
            }
        }

        // Get field value if possible
        var value interface{}
        if fieldValue.CanInterface() {
            value = fieldValue.Interface()
        }

        fieldInfo := FieldInfo{
            Name:       field.Name,
            Type:       field.Type.String(),
            Tags:       tags,
            Value:      value,
            IsExported: field.PkgPath == "", // Field is exported if PkgPath is empty
        }

        info.Fields = append(info.Fields, fieldInfo)
    }

    return info, nil
}

// PrettyPrint returns a formatted string representation of struct info
func (i *Inspector) PrettyPrint(info *StructInfo) string {
    var builder strings.Builder

    builder.WriteString(fmt.Sprintf("Struct: %s\n", info.Name))
    builder.WriteString("Fields:\n")

    for _, field := range info.Fields {
        builder.WriteString(fmt.Sprintf("  - %s:\n", field.Name))
        builder.WriteString(fmt.Sprintf("    Type: %s\n", field.Type))
        builder.WriteString(fmt.Sprintf("    Exported: %v\n", field.IsExported))

        if len(field.Tags) > 0 {
            builder.WriteString("    Tags:\n")
            for key, value := range field.Tags {
                builder.WriteString(fmt.Sprintf("      %s: %s\n", key, value))
            }
        }

        if field.IsExported && field.Value != nil {
            builder.WriteString(fmt.Sprintf("    Value: %v\n", field.Value))
        }
    }

    return builder.String()
}