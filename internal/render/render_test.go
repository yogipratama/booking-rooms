package render

import (
	"net/http"
	"testing"

	"github.com/yogipratama/booking-rooms/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var tmplData models.TemplateData
	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(request.Context(), "flash", "123")
	result := AddDefaultData(&tmplData, request)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTmpl(t *testing.T) {
	pathToTemplates = "./../../templates"
	tmplCache, err := CreateTmplCache()
	if err != nil {
		t.Error(err)
	}

	appConfig.TemplateCache = tmplCache

	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var writer myWriter

	err = RenderTmpl(&writer, request, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser")
	}

	err = RenderTmpl(&writer, request, "non-existent.page.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template but doesn't exist")
	}
}

func getSession() (*http.Request, error) {
	request, err := http.NewRequest("GET", "/whatever", nil)
	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	ctx, _ = session.Load(ctx, request.Header.Get("X-SESSION"))

	request = request.WithContext(ctx)

	return request, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(appConfig)
}

func TestCreateTmplCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTmplCache()
	if err != nil {
		t.Error(err)
	}
}
