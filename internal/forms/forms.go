package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom Form struct, embeds a URL values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a new Form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid checks if there are any errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks if the require values are passed in the form
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s cannot be blank.", field))
		}
	}
}

// Has checks if the form has a field
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength check if a field is of a minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	val := r.Form.Get(field)
	if len(val) < length {
		f.Errors.Add(field, fmt.Sprintf("%s must be at least %d characters long", field, length))
		return false
	}
	return true
}

// IsEmail check if input is a valid email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
