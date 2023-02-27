package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field can't be blank!")
		}
	}
}

// Has checks if form field is in post and not empty
func (form *Form) Has(field string, request *http.Request) bool {
	x := request.Form.Get(field)
	if x == "" {
		form.Errors.Add(field, "This field can't be blank!")
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (form *Form) MinLength(field string, length int, request *http.Request) bool {
	x := request.Form.Get(field)
	if len(x) < length {
		form.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail Checks for valid email address
func (form *Form) IsEmail(field string) {
	if !govalidator.IsEmail(form.Get(field)) {
		form.Errors.Add(field, "Invalid email address")
	}
}
