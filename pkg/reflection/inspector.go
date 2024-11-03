package reflection

import (
    "fmt"
    "reflect"
    "strings"

    "di-example/pkg/logger"
    "go.uber.org/zap"
)

type StructInfo struct {
    Name   string
    Fields []FieldInfo
}

type FieldInfo struct {
    Name       string
    Type       string
    Tags       map[string]string
    Value      interface{}
    IsExported bool
}

type Inspector struct {
    log *zap.SugaredLogger
}

func NewInspector() *Inspector {
    return &Inspector{
        log: logger.Get(),
    }
}

func (i *Inspector) InspectStruct(target interface{}) (*StructInfo, error) {
    i.log.Info("Starting struct inspection")

    // Check for nil target
    if target == nil {
        i.log.Error("Target is nil")
        return nil, fmt.Errorf("target cannot be nil")
    }

    targetValue := reflect.ValueOf(target)
    targetType := targetValue.Type()

    i.log.Debugw("Analyzing target type",
        "initialType", targetType.String())

    // Handle pointer types
    if targetType.Kind() == reflect.Ptr {
        i.log.Debug("Target is a pointer, dereferencing")
        // Check if pointer is nil
        if targetValue.IsNil() {
            i.log.Error("Target pointer is nil")
            return nil, fmt.Errorf("target pointer cannot be nil")
        }
        targetValue = targetValue.Elem()
        targetType = targetValue.Type()
    }

    // Ensure we're dealing with a struct
    if targetType.Kind() != reflect.Struct {
        i.log.Errorw("Invalid target type",
            "expectedKind", "struct",
            "actualKind", targetType.Kind())
        return nil, fmt.Errorf("target must be a struct or pointer to struct, got: %v",
            targetType.Kind())
    }

    i.log.Infow("Creating struct info",
        "structName", targetType.Name(),
        "numFields", targetType.NumField())

    info := &StructInfo{
        Name:   targetType.Name(),
        Fields: make([]FieldInfo, 0, targetType.NumField()),
    }

    // Analyze each field
    for fieldIdx := 0; fieldIdx < targetType.NumField(); fieldIdx++ {
        field := targetType.Field(fieldIdx)
        fieldValue := targetValue.Field(fieldIdx)

        i.log.Debugw("Analyzing field",
            "fieldName", field.Name,
            "fieldType", field.Type.String())

        // Parse field tags
        tags := make(map[string]string)
        if field.Tag != "" {
            i.log.Debugw("Parsing field tags",
                "fieldName", field.Name,
                "rawTags", field.Tag)

            for _, tag := range strings.Split(string(field.Tag), " ") {
                parts := strings.Split(tag, ":")
                if len(parts) == 2 {
                    key := strings.Trim(parts[0], "`")
                    value := strings.Trim(parts[1], `"`)
                    tags[key] = value
                }
            }
        }

        // Get field value if possible
        var value interface{}
        if fieldValue.CanInterface() {
            value = fieldValue.Interface()
            i.log.Debugw("Retrieved field value",
                "fieldName", field.Name,
                "canInterface", true)
        } else {
            i.log.Debugw("Cannot interface field value",
                "fieldName", field.Name,
                "canInterface", false)
        }

        isExported := field.PkgPath == ""
        i.log.Debugw("Field export status",
            "fieldName", field.Name,
            "isExported", isExported)

        fieldInfo := FieldInfo{
            Name:       field.Name,
            Type:       field.Type.String(),
            Tags:       tags,
            Value:      value,
            IsExported: isExported,
        }

        info.Fields = append(info.Fields, fieldInfo)
    }

    i.log.Info("Completed struct inspection")
    return info, nil
}

func (i *Inspector) PrettyPrint(info *StructInfo) string {
    i.log.Info("Generating pretty print output")

    var builder strings.Builder

    builder.WriteString(fmt.Sprintf("Struct: %s\n", info.Name))
    builder.WriteString("Fields:\n")

    for _, field := range info.Fields {
        i.log.Debugw("Pretty printing field",
            "fieldName", field.Name)

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