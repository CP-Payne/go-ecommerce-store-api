package user

import (
	"errors"
	"regexp"
)

type Email string

func ValidateEmail(e string) (Email, error) {
	match, _ := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", e)
	if !match {
		return "", errors.New("invalid email provided")
	}

	return Email(e), nil
}

type Name string

func ValidateName(n string) (Name, error) {
	if n == "" {
		return "", nil
	}
	match, _ := regexp.MatchString("^[A-Za-z ]+$", n)
	if !match {
		return "", errors.New("invalid name provided")
	}

	return Name(n), nil
}

type Password string

func ValidatePassword(p string) (Password, error) {
	// Check if it has at least one letter
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(p)
	// Check if it has at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(p)
	// Check if it has at least one special character
	hasSpecial := regexp.MustCompile(`[@$!%*#?&]`).MatchString(p)
	// Check if it has at least 8 characters
	isValidLength := len(p) >= 8

	if !hasLetter || !hasDigit || !hasSpecial || !isValidLength {
		return "", errors.New("password must be at least 8 characters long, contain at least one letter, one number, and one special character")
	}

	return Password(p), nil
}
