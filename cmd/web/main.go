package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/yogipratama/booking-rooms/internal/config"
	"github.com/yogipratama/booking-rooms/internal/handlers"
	"github.com/yogipratama/booking-rooms/internal/models"
	"github.com/yogipratama/booking-rooms/internal/render"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting app on port 8080")
	serve := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err = serve.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func run() error {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tmplCache, err := render.CreateTmplCache()
	if err != nil {
		log.Fatal("Can't create template cache")
		return err
	}

	app.TemplateCache = tmplCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return nil
}
