package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yogipratama/booking-rooms/internal/config"
	"github.com/yogipratama/booking-rooms/internal/models"
	"github.com/yogipratama/booking-rooms/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Creates a new Repository
func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

// sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	remoteIP := request.RemoteAddr
	repo.App.Session.Put(request.Context(), "remote_ip", remoteIP)

	render.RenderTmpl(writer, request, "home.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) About(writer http.ResponseWriter, request *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, im here!"

	remoteIP := repo.App.Session.GetString(request.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTmpl(writer, request, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "generals.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "majors.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) Availability(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "search-availability.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(writer http.ResponseWriter, request *http.Request) {
	start := request.Form.Get("start")
	end := request.Form.Get("end")
	writer.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) AvailabilityJSON(writer http.ResponseWriter, request *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	result, err := json.MarshalIndent(resp, "", "      ")
	if err != nil {
		log.Println(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(result)
}

func (repo *Repository) Contact(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "contact.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) Reservation(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "make-reservation.page.gohtml", &models.TemplateData{})
}
