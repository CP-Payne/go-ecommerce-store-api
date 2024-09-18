package user

import (
	"errors"
	"testing"
)

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		input    string
		expected Email
		err      error
	}{
		{"test@example.com", "test@example.com", nil},
		{"valid.email@domain.com", "valid.email@domain.com", nil},
		{"email@sub.domain.com", "email@sub.domain.com", nil},
		{"invalid-email", "", errors.New("invalid email provided")},
		{"missing@domain", "", errors.New("invalid email provided")},
		{"missingat.com", "", errors.New("invalid email provided")},
		{"", "", errors.New("invalid email provided")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ValidateEmail(tt.input)

			if err != nil && tt.err == nil {
				t.Fatalf("expected no error but got %v", err)
			}

			if err == nil && tt.err != nil {
				t.Fatalf("expected error %v but got no error", tt.err)
			}

			// If error is expected, check if it matches
			if err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("expected error %v but got %v", tt.err, err)
			}

			if result != tt.expected {
				t.Fatalf("expected result %v but got %v", tt.expected, result)
			}
		})
	}
}

func TestNameValidation(t *testing.T) {
	tests := []struct {
		input    string
		expected Name
		err      error
	}{
		{"John", "John", nil},
		{"John Doe", "John Doe", nil},
		{"Invalid Name321", "", errors.New("invalid name provided")},
		{"23432", "", errors.New("invalid name provided")},
		{"Invalid_Name", "", errors.New("invalid name provided")},
		{"", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ValidateName(tt.input)

			if err != nil && tt.err == nil {
				t.Fatalf("expected no error but got %v", err)
			}

			if err == nil && tt.err != nil {
				t.Fatalf("expected error %v but got no error", tt.err)
			}

			// If error is expected, check if it matches
			if err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("expected error %v but got %v", tt.err, err)
			}

			if result != tt.expected {
				t.Fatalf("expected result %v but got %v", tt.expected, result)
			}
		})
	}
}

func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		input    string
		expected Password
		err      error
	}{
		{"password1!", "password1!", nil},
		{"12345678A#", "12345678A#", nil},
		{"pass", "", errors.New("password must be at least 8 characters long, contain at least one letter, one number, and one special character")},
		{"password", "", errors.New("password must be at least 8 characters long, contain at least one letter, one number, and one special character")},
		{"password1", "", errors.New("password must be at least 8 characters long, contain at least one letter, one number, and one special character")},
		{"password!", "", errors.New("password must be at least 8 characters long, contain at least one letter, one number, and one special character")},
		{"", "", errors.New("password must be at least 8 characters long, contain at least one letter, one number, and one special character")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ValidatePassword(tt.input)

			if err != nil && tt.err == nil {
				t.Fatalf("expected no error but got %v", err)
			}

			if err == nil && tt.err != nil {
				t.Fatalf("expected error %v but got no error", tt.err)
			}

			// If error is expected, check if it matches
			if err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("expected error %v but got %v", tt.err, err)
			}

			if result != tt.expected {
				t.Fatalf("expected result %v but got %v", tt.expected, result)
			}
		})
	}
}
