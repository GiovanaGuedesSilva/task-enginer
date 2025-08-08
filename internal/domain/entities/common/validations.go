package common

import (
	"fmt"
	"time"
)

type FieldValidationError struct {
	Field   string
	Message string
}

func (e FieldValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type ValidationErrors []FieldValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}

	msg := "validation errors:"
	for _, err := range e {
		msg += fmt.Sprintf("\n- %s: %s", err.Field, err.Message)
	}
	return msg
}

// String validations

func ValidateRequired(fieldName, value string) *FieldValidationError {
	if value == "" {
		return &FieldValidationError{
			Field:   fieldName,
			Message: "field is required",
		}
	}
	return nil
}

func ValidateStringLength(fieldName, value string, maxLength int) *FieldValidationError {
	if len(value) > maxLength {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field cannot exceed %d characters", maxLength),
		}
	}
	return nil
}

func ValidateStringLengthRange(fieldName, value string, minLength, maxLength int) *FieldValidationError {
	if len(value) < minLength {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field must have at least %d characters", minLength),
		}
	}
	if len(value) > maxLength {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field cannot exceed %d characters", maxLength),
		}
	}
	return nil
}

func ValidatePositiveInt(fieldName string, value int64) *FieldValidationError {
	if value <= 0 {
		return &FieldValidationError{
			Field:   fieldName,
			Message: "field must be a positive integer",
		}
	}
	return nil
}

func ValidatePositiveFloat(fieldName string, value float64) *FieldValidationError {
	if value < 0 {
		return &FieldValidationError{
			Field:   fieldName,
			Message: "field cannot be negative",
		}
	}
	return nil
}

func ValidateEnum(fieldName string, value interface{}, allowedValues ...interface{}) *FieldValidationError {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return &FieldValidationError{
		Field:   fieldName,
		Message: fmt.Sprintf("field must be one of: %v", allowedValues),
	}
}

// Date validations

func ValidateDateRange(fieldName string, date time.Time, minDate, maxDate time.Time) *FieldValidationError {
	if date.Before(minDate) {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field cannot be before %s", minDate.Format("2006-01-02")),
		}
	}
	if date.After(maxDate) {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field cannot be after %s", maxDate.Format("2006-01-02")),
		}
	}
	return nil
}

func ValidateDateOrder(startFieldName, endFieldName string, startDate, endDate time.Time) *FieldValidationError {
	if !startDate.IsZero() && !endDate.IsZero() {
		if startDate.After(endDate) {
			return &FieldValidationError{
				Field:   fmt.Sprintf("%s and %s", startFieldName, endFieldName),
				Message: fmt.Sprintf("%s cannot be after %s", startFieldName, endFieldName),
			}
		}
	}
	return nil
}

// Email and URL validations

func ValidateEmail(fieldName, email string) *FieldValidationError {
	if email == "" {
		return nil
	}

	if len(email) < 5 || !contains(email, "@") || !contains(email, ".") {
		return &FieldValidationError{
			Field:   fieldName,
			Message: "field must be a valid email address",
		}
	}
	return nil
}

func ValidateURL(fieldName, url string) *FieldValidationError {
	if url == "" {
		return nil
	}

	if len(url) < 10 || (!contains(url, "http://") && !contains(url, "https://")) {
		return &FieldValidationError{
			Field:   fieldName,
			Message: "field must be a valid URL",
		}
	}
	return nil
}

// Numeric validations

func ValidateNumericRange(fieldName string, value, min, max float64) *FieldValidationError {
	if value < min {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field must be at least %.2f", min),
		}
	}
	if value > max {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field cannot exceed %.2f", max),
		}
	}
	return nil
}

func ValidateIntRange(fieldName string, value, min, max int64) *FieldValidationError {
	if value < min {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field must be at least %d", min),
		}
	}
	if value > max {
		return &FieldValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field cannot exceed %d", max),
		}
	}
	return nil
}

// Domain specific validations (Project)

func ValidateProjectStatus(fieldName string, status ProjectStatus) *FieldValidationError {
	return ValidateEnum(fieldName, status,
		ProjectStatusActive,
		ProjectStatusInactive,
		ProjectStatusArchived,
		ProjectStatusDeleted,
	)
}

func ValidateProjectPriority(fieldName string, priority ProjectPriority) *FieldValidationError {
	return ValidateEnum(fieldName, priority,
		ProjectPriorityLow,
		ProjectPriorityMedium,
		ProjectPriorityHigh,
		ProjectPriorityUrgent,
	)
}

func ValidateTaskStatus(fieldName string, status TaskStatus) *FieldValidationError {
	return ValidateEnum(fieldName, status,
		TaskStatusPending,
		TaskStatusInProgress,
		TaskStatusCompleted,
		TaskStatusCancelled,
	)
}

func ValidateTaskPriority(fieldName string, priority TaskPriority) *FieldValidationError {
	return ValidateEnum(fieldName, priority,
		TaskPriorityLow,
		TaskPriorityMedium,
		TaskPriorityHigh,
		TaskPriorityUrgent,
	)
}

func ValidateUserRole(fieldName string, role UserRole) *FieldValidationError {
	return ValidateEnum(fieldName, role,
		UserRoleAdmin,
		UserRoleManager,
		UserRoleMember,
		UserRoleObserver,
	)
}

func ValidateUserStatus(fieldName string, status UserStatus) *FieldValidationError {
	return ValidateEnum(fieldName, status,
		UserStatusActive,
		UserStatusInactive,
		UserStatusSuspended,
	)
}

// Auxiliary functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Convenience functions for multiple validations

func ValidateFields(validations ...*FieldValidationError) error {
	var errors ValidationErrors

	for _, validation := range validations {
		if validation != nil {
			errors = append(errors, *validation)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// ValidateRequiredFields for multiple required fields
func ValidateRequiredFields(fields map[string]string) error {
	var errors ValidationErrors

	for fieldName, value := range fields {
		if err := ValidateRequired(fieldName, value); err != nil {
			errors = append(errors, *err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// ValidateStringFields for multiple string fields with maximum length
func ValidateStringFields(fields map[string]struct {
	Value     string
	MaxLength int
}) error {
	var errors ValidationErrors

	for fieldName, field := range fields {
		if err := ValidateStringLength(fieldName, field.Value, field.MaxLength); err != nil {
			errors = append(errors, *err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
