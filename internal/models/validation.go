package models

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/Pomog/ForumFFF/internal/forms"
)

type Validatable interface {
	Validate() []error
}

// ValidationConfig holds configuration parameters for validation.
type ValidationConfig struct {
	MinLen           int
	MaxLen           int
	PasswordMinLen   int
	PasswordMaxLen   int
	MinCategoryLen   int
	MaxCategoryLen   int
	MinSubjectLen    int
	MaxSubjectLen    int
	SingleWordMaxLen int
}

/*
Validate performs validation on the User struct.
It checks for the presence of required fields and validates the email format.
Returns a slice of errors, where each error represents a validation issue.
*/
func (u *User) Validate(config ValidationConfig) []error {
	var validationErrors []error

	// Validate username
	if err := ValidateRequired("UserName", u.UserName); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Validate email
	if err := ValidateRequired("Email", u.Email); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if _, err := mail.ParseAddress(u.Email); err != nil {
			validationErrors = append(validationErrors, errors.New("invalid email address"))
		}
	}

	// Validate first name
	if err := ValidateRequired("FirstName", u.FirstName); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("FirstName", u.FirstName, config.MinLen, config.MaxLen); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	// Validate last name
	if err := ValidateRequired("LastName", u.LastName); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("LastName", u.LastName, config.MinLen, config.MaxLen); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	// Validate password
	if err := ValidateRequired("Password", u.Password); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("Password", u.Password, config.PasswordMinLen, config.PasswordMaxLen); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	// Add more validation logic for other fields

	return validationErrors
}

/*
Validate performs validation on the thread struct.
It checks for the presence of required fields and validates the length format.
Returns a slice of errors, where each error represents a validation issue.
*/
func (thread *Thread) Validate(config ValidationConfig) []error {
	var validationErrors []error

	// Validate Category
	if err := ValidateRequired("Category", thread.Category); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("Category", thread.Category, config.MinCategoryLen, config.MaxCategoryLen); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	// Validate Subject
	if err := ValidateRequired("Subject", thread.Subject); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("Subject", thread.Subject, config.MinSubjectLen, config.MaxSubjectLen); err != nil {
			validationErrors = append(validationErrors, err)
		} else {
			if !forms.CheckSingleWordLen(thread.Subject, config.SingleWordMaxLen) {
				err := fmt.Errorf("the Subject without spaces is not allowed, max len of each word (without spaces) is %d", config.SingleWordMaxLen)
				validationErrors = append(validationErrors, err)
			}
		}
	}

	// Add more validation logic for other fields

	return validationErrors
}

/*
Validate performs validation on the thread struct.
It checks for the presence of required fields and validates the length format.
Returns a slice of errors, where each error represents a validation issue.
*/
func (post *Post) Validate(config ValidationConfig) []error {
	var validationErrors []error

	// Validate Content
	if err := ValidateRequired("Content", post.Content); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("Content", post.Content, config.MinSubjectLen, config.MaxSubjectLen); err != nil {
			validationErrors = append(validationErrors, err)
		} else {
			if !forms.CheckSingleWordLen(post.Content, config.SingleWordMaxLen) {
				err := fmt.Errorf("the Content without spaces is not allowed, max len of each word (without spaces) is %d", config.SingleWordMaxLen)
				validationErrors = append(validationErrors, err)
			}
		}
	}

	// Add more validation logic for other fields

	return validationErrors
}

func (pm *PM) ValidatePM(config ValidationConfig) []error {
	var validationErrors []error

	// Validate Content
	if err := ValidateRequired("Content", pm.Content); err != nil {
		validationErrors = append(validationErrors, err)
	} else {
		if err := ValidateLength("Content", pm.Content, config.MinSubjectLen, config.MaxSubjectLen); err != nil {
			validationErrors = append(validationErrors, err)
		} else {
			if !forms.CheckSingleWordLen(pm.Content, config.SingleWordMaxLen) {
				err := fmt.Errorf("the Content without spaces is not allowed, max len of each word (without spaces) is %d", config.SingleWordMaxLen)
				validationErrors = append(validationErrors, err)
			}
		}
	}

	// Add more validation logic for other fields

	return validationErrors
}


// ValidateRequired checks if the field is required and not empty.
func ValidateRequired(fieldName string, fieldValue string) error {
	if fieldValue == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}

// ValidateLength checks if the length of the field is within the specified range.
func ValidateLength(fieldName string, fieldValue string, minLen int, maxLen int) error {
	length := len(strings.TrimSpace(fieldValue))
	if length < minLen || length > maxLen {
		return errors.New(fieldName + " length must be between " + fmt.Sprint(minLen) + " and " + fmt.Sprint(maxLen))
	}
	return nil
}
