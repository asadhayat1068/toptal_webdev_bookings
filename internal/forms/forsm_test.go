package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/endpoint", nil)
	f := New(r.PostForm)
	if !f.Valid() {
		t.Error("Invalidate a valid form")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/endpoint", nil)
	f := New(r.PostForm)
	f.Required("first_name")
	if f.Valid() {
		t.Error("Validate an invalid form")
	}

	data := url.Values{}
	data.Add("first_name", "First")
	r.PostForm = data
	f = New(r.PostForm)
	if !f.Valid() {
		t.Error("Invalidate a valid form")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/endpoint", nil)
	f := New(r.PostForm)
	if f.Has("last_name") {
		t.Error("Found a value that is not present in form")
	}

	data := url.Values{}
	data.Add("first_name", "First")
	r.PostForm = data
	f = New(r.PostForm)
	if !f.Has("first_name") {
		t.Error("Does not found a value that is present in form")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/endpoint", nil)
	data := url.Values{}
	data.Add("first_name", "First")
	r.PostForm = data
	f := New(r.PostForm)
	if f.MinLength("first_name", 6) {
		t.Error("Min-length returned true when should be false")
	}
	if !f.MinLength("first_name", 5) {
		t.Error("Min-length returned false when should be true")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/endpoint", nil)
	data := url.Values{}
	data.Add("invalid_email", "NotAnEmail")
	data.Add("valid_email", "valid@email.com")
	r.PostForm = data
	f := New(r.PostForm)
	if f.IsEmail("invalid_email") {
		t.Error("non-email value is validated as email.")
	}
	if !f.IsEmail("valid_email") {
		t.Error("valid-email is returned as invalid")
	}
}
