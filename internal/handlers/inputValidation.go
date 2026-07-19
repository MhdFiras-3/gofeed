package handlers

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

type inputError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func validateInput(name, email, password string) []inputError {
	var errs []inputError

	if err := parseName(name); err != nil {
		errs = append(errs, inputError{Field: "name", Message: err.Error()})
	}
	if err := parseEmail(email); err != nil {
		errs = append(errs, inputError{Field: "email", Message: err.Error()})
	}
	if err := parsePassword(password); err != nil {
		errs = append(errs, inputError{Field: "password", Message: err.Error()})
	}
	return errs
}

func parseEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty!")
	}
	if len(email) > 50 {
		return errors.New("email too long")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address %w", err)
	}
	return nil
}

func parsePassword(password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) > 100 {
		return errors.New("password too long")
	}
	if len(password) < 20 {
		return errors.New("password too short")
	}

	return nil

}

func parseName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) > 30 {
		return errors.New("name too long")
	}
	if containsControlChars(name) {
		return errors.New("name contains invalid characters")
	}
	return nil
}

func containsControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
