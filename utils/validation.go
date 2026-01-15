package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) []string {
	var errors []string

	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}

	if len(password) > 128 {
		errors = append(errors, "Password must be less than 128 characters")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}

	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}

	if !hasNumber {
		errors = append(errors, "Password must contain at least one number")
	}

	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}

	return errors
}

// ValidateProductName validates product name
func ValidateProductName(name string) []string {
	var errors []string

	if strings.TrimSpace(name) == "" {
		errors = append(errors, "Product name cannot be empty")
	}

	if len(name) < 3 {
		errors = append(errors, "Product name must be at least 3 characters long")
	}

	if len(name) > 100 {
		errors = append(errors, "Product name must be less than 100 characters")
	}

	// Check for potentially dangerous characters
	dangerousChars := []string{"<", ">", "&", "\"", "'", "/", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(name, char) {
			errors = append(errors, fmt.Sprintf("Product name cannot contain '%s'", char))
		}
	}

	return errors
}

// ValidatePrice validates price
func ValidatePrice(price float64) []string {
	var errors []string

	if price <= 0 {
		errors = append(errors, "Price must be greater than 0")
	}

	if price > 999999.99 {
		errors = append(errors, "Price cannot exceed 999,999.99")
	}

	return errors
}

// ValidateStock validates stock quantity
func ValidateStock(stock uint) []string {
	var errors []string

	if stock > 1000000 {
		errors = append(errors, "Stock cannot exceed 1,000,000 units")
	}

	return errors
}

// HandleValidationError formats validation errors for API responses
func HandleValidationError(c *gin.Context, errors []string) {
	if len(errors) > 0 {
		c.JSON(400, gin.H{
			"error":   "Validation failed",
			"details": errors,
		})
	}
}
