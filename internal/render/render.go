package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/yogipratama/booking-rooms/internal/config"
	"github.com/yogipratama/booking-rooms/internal/models"
)

var appConfig *config.AppConfig

// sets the config for the template package
func NewTemplates(app *config.AppConfig) {
	appConfig = app
}

func AddDefaultData(tmplData *models.TemplateData, request *http.Request) *models.TemplateData {
	tmplData.Flash = appConfig.Session.PopString(request.Context(), "flash")
	tmplData.Warning = appConfig.Session.PopString(request.Context(), "warning")
	tmplData.Error = appConfig.Session.PopString(request.Context(), "error")
	tmplData.CSRFToken = nosurf.Token(request)
	return tmplData
}

func RenderTmpl(writer http.ResponseWriter, request *http.Request, tmpl string, tmplData *models.TemplateData) {

	var tmplCache map[string]*template.Template
	if appConfig.UseCache {
		// get the template cache from the app config
		tmplCache = appConfig.TemplateCache
	} else {
		tmplCache, _ = CreateTmplCache()
	}

	// get requested template from cache
	template, ok := tmplCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, request)

	err := template.Execute(buf, tmplData)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(writer)
	if err != nil {
		log.Println(err)
	}
}

func CreateTmplCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	// get all of the files named *.gohtml from ./templates
	files, err := filepath.Glob("./templates/*.gohtml")
	if err != nil {
		return tmplCache, err
	}

	// range through all files ending with *.gohtml
	for _, page := range files {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return tmplCache, err
		}

		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = tmplSet
	}

	return tmplCache, nil
}
