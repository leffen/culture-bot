package tpl

import (
	"html/template"

	"github.com/gocraft/web"
	log "github.com/sirupsen/logrus"
)

// Data struct
type Data struct {
	TemplateFile string
	Data         interface{}
}

// Render func
func (s Data) Render(w web.ResponseWriter, r *web.Request) {
	defer r.Body.Close()

	if s.TemplateFile == "" {
		log.Error("empty TemplateFile")
		return
	}

	t, err := template.ParseFiles("dashboard/templates/"+s.TemplateFile, "dashboard/templates/layout.html")
	if err != nil {
		log.WithError(err).Error("unable to parse files")
		return
	}

	err = t.ExecuteTemplate(w, "base", s)
	if err != nil {
		log.WithError(err).Error("unable to execute template")
	}
}
