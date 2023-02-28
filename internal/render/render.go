package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/yogipratama/booking-rooms/internal/config"
	"github.com/yogipratama/booking-rooms/internal/models"
)

var appConfig *config.AppConfig

var pathToTemplates = "./templates"
var functions = template.FuncMap{}

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

func RenderTmpl(writer http.ResponseWriter, request *http.Request, tmpl string, tmplData *models.TemplateData) error {

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
		return errors.New("can't get template from cache")
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
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

func CreateTmplCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	// get all of the files named *.gohtml from ./templates
	files, err := filepath.Glob(fmt.Sprintf("%s/*.gohtml", pathToTemplates))
	if err != nil {
		return tmplCache, err
	}

	// range through all files ending with *.gohtml
	for _, page := range files {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return tmplCache, err
		}

		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = tmplSet
	}

	return tmplCache, nil
}
