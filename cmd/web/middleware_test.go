package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var handler myHandler
	h := NoSurf(&handler)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("Type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var handler myHandler
	h := SessionLoad(&handler)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("Type is not http.Handler, but is %T", v))
	}
}
