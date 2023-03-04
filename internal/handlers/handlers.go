package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yogipratama/booking-rooms/internal/config"
	"github.com/yogipratama/booking-rooms/internal/driver"
	"github.com/yogipratama/booking-rooms/internal/forms"
	"github.com/yogipratama/booking-rooms/internal/helpers"
	"github.com/yogipratama/booking-rooms/internal/models"
	"github.com/yogipratama/booking-rooms/internal/render"
	"github.com/yogipratama/booking-rooms/internal/repository"
	"github.com/yogipratama/booking-rooms/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates a new Repository
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}

// sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "home.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) About(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "about.page.gohtml", &models.TemplateData{})
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

func (repo *Repository) Contact(writer http.ResponseWriter, request *http.Request) {
	render.RenderTmpl(writer, request, "contact.page.gohtml", &models.TemplateData{})
}

func (repo *Repository) Reservation(writer http.ResponseWriter, request *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTmpl(writer, request, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (repo *Repository) PostReservation(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		helpers.ServerError(writer, err)
	}

	reservation := models.Reservation{
		FirstName: request.Form.Get("first_name"),
		LastName:  request.Form.Get("last_name"),
		Email:     request.Form.Get("email"),
		Phone:     request.Form.Get("phone"),
	}

	form := forms.New(request.PostForm)

	// form.Has("first_name", request)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 4)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTmpl(writer, request, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	repo.App.Session.Put(request.Context(), "reservation", reservation)

	http.Redirect(writer, request, "/reservation-summary", http.StatusSeeOther)
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
		helpers.ServerError(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(result)
}

func (repo *Repository) ReservationSummary(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := repo.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.ErrorLog.Println("Can't get error from session")
		repo.App.Session.Put(request.Context(), "error", "Cannot get reservation from session")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.App.Session.Remove(request.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTmpl(writer, request, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}
