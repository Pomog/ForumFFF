package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Pomog/ForumFFF/pkg/config"
	"github.com/Pomog/ForumFFF/pkg/models"
)

// this var serves to pass data from main.go to render.go
var app *config.AppConfig

var pathToTemplates = "./template"
var functions = template.FuncMap{}

// NewTemplate sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RendererTemplate renders template using html/template
func RendererTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		//get the template cache from AppConfig8
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
		// log.Println("Using CreateTemplateCache")
	}

	//get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	td = AddDefaultData(td)

	//optional final error check
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		// log.Println("error writting template to brwoser ", err)
	}

}

// create a template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all files *.page.tmpl from templates ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		// log.Println("error getting pages from templates ", err)
		return myCache, err
	}
	// log.Println("pages: ", pages)

	// range over pages
	for _, page := range pages {
		// get file name
		name := filepath.Base(page)
		// log.Println("currently parsing page: ", page)
		// parse page
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			// log.Println("error parsing page ", err)
			return myCache, err
		}
		// get base layout
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			// log.Println("error getting base layout ", err)
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				// log.Println("error parsing base layout ", err)
				return myCache, err
			}
		}
		// add to cache
		myCache[name] = ts
	}
	return myCache, nil
}
