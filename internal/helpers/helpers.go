package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/yogipratama/booking-rooms/internal/config"
)

var appConfig *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(app *config.AppConfig) {
	appConfig = app
}

func ClientError(writer http.ResponseWriter, statusCode int) {
	appConfig.InfoLog.Println("Client error with status of", statusCode)
	http.Error(writer, http.StatusText(statusCode), statusCode)
}

func ServerError(writer http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	appConfig.ErrorLog.Println(trace)

	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
