package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/models"
)

// this var serves to pass data from main.go to render.go
var app *config.AppConfig

var pathToTemplates = "./template"
var functions = template.FuncMap{
	"postsLen": func(allPosts []models.Post) int {
		return len(allPosts)
	},
	"findLastPost": func(allPosts []models.Post) models.Post {
		var latestPost models.Post
		latestPost.Created, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		for _, post := range allPosts {
			if post.Created.After(latestPost.Created) {
				latestPost = post
			}
		}
		return latestPost
	},
	"numberOfPosts": func(allPosts []models.Post) int {
		return len(allPosts)
	},
	"convertTime": func(post models.Post) string {
		return post.Created.Format("2006-01-02 15:04:05")
	},
	"shortenPost": func(allPosts []models.Post) string {
		var latestPost2 models.Post
		latestPost2.Created, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		for _, post := range allPosts {
			if post.Created.After(latestPost2.Created) {
				latestPost2 = post
			}
		}
		pst := latestPost2.Content
		if len(pst) <= 80 {
			return pst
		}
		return pst[0:80]

	},
}

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

// CreateTemplateCache create a template cache
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
