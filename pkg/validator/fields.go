package validator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var parts []string
	for _, err := range e {
		if err.Field != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", err.Field, err.Message))
		} else {
			parts = append(parts, err.Message)
		}
	}
	return strings.Join(parts, "; ")
}

type ValidationRule func(value interface{}) *ValidationError

func ValidateField(fieldName string, value interface{}, rules ...ValidationRule) []ValidationError {
	var errs []ValidationError
	for _, rule := range rules {
		if err := rule(value); err != nil {
			err.Field = fieldName
			errs = append(errs, *err)
		}
	}
	return errs
}

func RequiredString() ValidationRule {
	return func(value interface{}) *ValidationError {
		str, ok := value.(string)
		if !ok {
			return &ValidationError{Message: "invalid type, expected string"}
		}
		if strings.TrimSpace(str) == "" {
			return &ValidationError{Message: "cannot be empty"}
		}
		return nil
	}
}

func MaxLength(max int) ValidationRule {
	return func(value interface{}) *ValidationError {
		str, ok := value.(string)
		if !ok {
			return &ValidationError{Message: "invalid type, expected string"}
		}
		if len(str) > max {
			return &ValidationError{Message: fmt.Sprintf("cannot exceed %d characters", max)}
		}
		return nil
	}
}

func MinInt(min int64) ValidationRule {
	return func(value interface{}) *ValidationError {
		v, ok := value.(int64)
		if !ok {
			return &ValidationError{Message: "invalid type, expected int64"}
		}
		if v < min {
			return &ValidationError{Message: fmt.Sprintf("must be at least %d", min)}
		}
		return nil
	}
}

func MinFloat(min float64) ValidationRule {
	return func(value interface{}) *ValidationError {
		v, ok := value.(float64)
		if !ok {
			return &ValidationError{Message: "invalid type, expected float64"}
		}
		if v < min {
			return &ValidationError{Message: fmt.Sprintf("must be at least %.2f", min)}
		}
		return nil
	}
}

func ValidateEnum(isValidFunc func(interface{}) bool) ValidationRule {
	return func(value interface{}) *ValidationError {
		if !isValidFunc(value) {
			return &ValidationError{Message: "invalid value"}
		}
		return nil
	}
}
