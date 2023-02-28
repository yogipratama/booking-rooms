package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yogipratama/booking-rooms/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("Type is not http.Handler, but is %T", v))
	}
}
