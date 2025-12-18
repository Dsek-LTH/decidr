package templates

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed **/*.tmpl
var tmplFS embed.FS

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("dict requires an even number of arguments")
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m, nil
}

func NewTemplateRenderer() *TemplateRenderer {
	funcMap := template.FuncMap{
		"dict": dict,
	}
	templates := make(map[string]*template.Template)

	base, err := template.New("base").Funcs(funcMap).
		ParseFS(tmplFS, "layouts/*.tmpl", "components/*.tmpl")
	if err != nil {
		panic("failed to parse base layouts and partials: " + err.Error())
	}

	pageFiles, err := fs.Glob(tmplFS, "pages/*.tmpl")
	if err != nil {
		panic("failed to glob page template files: " + err.Error())
	}

	for _, pageFile := range pageFiles {
		tmplName := strings.TrimSuffix(filepath.Base(pageFile), ".tmpl")

		baseCopy, err := base.Clone()
		if err != nil {
			panic("failed to clone base template for " + tmplName + ": " + err.Error())
		}

		tmpl, err := baseCopy.ParseFS(tmplFS, pageFile)
		if err != nil {
			panic("failed to parse page template " + pageFile + ": " + err.Error())
		}

		templates[tmplName] = tmpl
	}

	return &TemplateRenderer{templates: templates}
}

// Render writes the template output into w.
// The template identifier may include the block separator (blockName = page#block).
func (tr *TemplateRenderer) Render(w io.Writer, templateIdentifier string, data any) error {
	pageName, blockName, hasBlock := strings.Cut(templateIdentifier, "#")
	if !hasBlock {
		blockName = pageName
	}
	if blockName == "" {
		return fmt.Errorf("error rendering template '%s': block templateIdentifier cannot be empty", templateIdentifier)
	}

	tmpl, ok := tr.templates[pageName]
	if !ok {
		return fmt.Errorf("template with pageName '%s' not found", pageName)
	}

	return tmpl.ExecuteTemplate(w, blockName, data)
}
