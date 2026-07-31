package handlers

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strconv"
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
	if len(name) > 45 {
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
func parseURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return errors.New("url cannot be empty")
	}
	if len(rawURL) > 150 {
		return errors.New("url too long")
	}
	url, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalde url: %w", err)
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return errors.New("url must use http or https")
	}
	return nil
}

func parseQueryParamLimit(limit string) int32 {
	if limit == "" {
		return 10
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 10
	}
	if limitInt <= 0 {
		return 10
	}
	if limitInt > 40 {
		return 40
	}
	return int32(limitInt)
}
