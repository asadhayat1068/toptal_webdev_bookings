package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not HttpHandler, but is %T", v)

	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	s := SessionLoad(&myH)
	switch v := s.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not HttpHandler, but is %T", v)

	}
}
