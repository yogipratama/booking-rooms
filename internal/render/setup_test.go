package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/yogipratama/booking-rooms/internal/config"
	"github.com/yogipratama/booking-rooms/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	appConfig = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tmplWriter *myWriter) Header() http.Header {
	var header http.Header
	return header
}

func (tmplWriter *myWriter) WriteHeader(code int) {

}

func (tmplWriter *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
