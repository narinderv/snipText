package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Structure to hold Form field data and errors
type FormInfo struct {
	url.Values
	Errors errors
}

// Initialize a new form
func NewForm(formData url.Values) *FormInfo {
	return &FormInfo{
		formData,
		errors(map[string][]string{}),
	}
}

// Form Validation functions
// Check if form is valid
func (form *FormInfo) IsValid() bool {
	return len(form.Errors) == 0
}

// Check Required fields
func (form *FormInfo) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// Check Max length
func (form *FormInfo) MaxLength(field string, maxLen int) {
	value := form.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > maxLen {
		form.Errors.Add(field, fmt.Sprintf("This field is too long (max. length: %d)", maxLen))
	}
}

// Check Max length
func (form *FormInfo) MinLength(field string, minLen int) {
	value := form.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < minLen {
		form.Errors.Add(field, fmt.Sprintf("This field is too short (min. length: %d)", minLen))
	}
}

// Check Max length
func (form *FormInfo) MatchPattern(field string, pattern *regexp.Regexp) {
	value := form.Get(field)

	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		form.Errors.Add(field, "This field is invalid")
	}
}

// Check permitted values
func (form *FormInfo) Permittedvalues(field string, validVals ...string) {
	value := form.Get(field)

	if value == "" {
		return
	}

	for _, validVal := range validVals {
		if value == validVal {
			return
		}
	}

	form.Errors.Add(field, "This field is invalid")
}
