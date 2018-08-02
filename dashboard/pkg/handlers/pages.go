package handlers

import (
	"net/http"

	"github.com/plutov/culture-bot/dashboard/pkg/greeter"
	"github.com/plutov/culture-bot/dashboard/pkg/recognition"
	"github.com/plutov/culture-bot/dashboard/pkg/stats"
	"github.com/plutov/culture-bot/dashboard/pkg/tpl"

	"github.com/gocraft/web"
)

func (c *Context) home(w web.ResponseWriter, r *web.Request) {
	d := tpl.Data{
		TemplateFile: "home.html",
		Data: struct {
			Error string
		}{},
	}

	d.Render(w, r)
}

func (c *Context) client(w web.ResponseWriter, r *web.Request) {
	greeter.ClearGreetings()

	d := tpl.Data{
		TemplateFile: "client.html",
		Data: struct {
			Error string
		}{},
	}

	d.Render(w, r)
}

func (c *Context) people(w web.ResponseWriter, r *web.Request) {
	users := stats.GetUsers()
	for _, u := range users {
		u.AverageRecognition = stats.GetAverageRecognitionValue(u.Slack)
		userPhotos, err := recognition.GetUserPhotos(u.Slack)
		if err == nil {
			for _, p := range userPhotos {
				u.Photos = append(u.Photos, p.Filename)
			}
		}
	}

	d := tpl.Data{
		TemplateFile: "people.html",
		Data: struct {
			Error string
			Users []*stats.User
		}{
			Users: users,
		},
	}

	d.Render(w, r)
}

func (c *Context) settings(w web.ResponseWriter, r *web.Request) {
	d := tpl.Data{
		TemplateFile: "settings.html",
		Data: struct {
			Error string
		}{},
	}

	d.Render(w, r)
}

func (c *Context) deleteUser(w web.ResponseWriter, r *web.Request) {
	id := r.URL.Query().Get("id")

	if len(id) > 0 {
		stats.DeleteUser(id)
	}

	http.Redirect(w, r.Request, "/people", http.StatusFound)
}

func (c *Context) updateUser(w web.ResponseWriter, r *web.Request) {
	id := r.FormValue("id")

	if len(id) > 0 {
		stats.UpdateUser(id, stats.User{
			First:     r.FormValue("first"),
			Last:      r.FormValue("last"),
			Slack:     id,
			Pronounce: r.FormValue("pronounce"),
		})
	}

	http.Redirect(w, r.Request, "/profile?id="+id, http.StatusFound)
}

func (c *Context) addUser(w web.ResponseWriter, r *web.Request) {
	stats.AddUser(stats.User{
		First:     r.FormValue("first"),
		Last:      r.FormValue("last"),
		Slack:     r.FormValue("slack"),
		Pronounce: r.FormValue("pronounce"),
	})

	http.Redirect(w, r.Request, "/profile?id="+r.FormValue("slack"), http.StatusFound)
}
