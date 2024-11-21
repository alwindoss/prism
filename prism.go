package prism

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func (p *prismRender) printTemplateCache() {
	for n, v := range p.templateCache {
		fmt.Printf("[%s]\t=\t%s\n", n, v.Name())
	}
}

// createTemplateCache initializes the template cache
func (p *prismRender) createTemplateCache() {
	tc := make(map[string]*template.Template)

	// Load all layouts, pages, and partials
	layouts, err := filepath.Glob(p.layoutPath)
	if err != nil {
		panic(err)
	}
	pages, err := filepath.Glob(p.pagesPath)
	if err != nil {
		panic(err)
	}
	partials, err := filepath.Glob(p.partialsPath)
	if err != nil {
		panic(err)
	}

	// Combine templates
	for _, page := range pages {
		files := append(layouts, page)
		files = append(files, partials...)
		templateName := filepath.Base(page)
		tc[templateName] = template.Must(template.ParseFiles(files...))
	}
	p.templateCache = tc
}

type Renderer interface {
	Render(w http.ResponseWriter, templateName string, layoutName string, data interface{})
}

type Config struct {
	// the path should resemble templates/layouts/*.html
	LayoutPath string

	// the path should resemble templates/pages/*.html
	PagesPath string

	// the path should resemble templates/partials/*.html
	PartialsPath string
}

func New(cfg *Config) Renderer {
	pr := &prismRender{
		layoutPath:   cfg.LayoutPath,
		pagesPath:    cfg.PagesPath,
		partialsPath: cfg.PartialsPath,
	}
	pr.createTemplateCache()
	pr.printTemplateCache()
	return pr
}

type prismRender struct {
	templateCache map[string]*template.Template
	layoutPath    string
	pagesPath     string
	partialsPath  string
}

// Render implements Renderer.
func (p *prismRender) Render(w http.ResponseWriter, templateName string, layoutName string, data interface{}) {
	tmpl, exists := p.templateCache[templateName]
	if !exists {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err := tmpl.ExecuteTemplate(w, layoutName, data)
	if err != nil {
		fmt.Println("Error: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
