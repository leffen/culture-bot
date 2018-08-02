package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gocraft/web"
	log "github.com/sirupsen/logrus"
	"github.com/plutov/culture-bot/dashboard/pkg/recognition"
	"github.com/plutov/culture-bot/dashboard/pkg/stats"
	"github.com/plutov/culture-bot/dashboard/pkg/tpl"
)

func (c *Context) newProfile(w web.ResponseWriter, r *web.Request) {
	d := tpl.Data{
		TemplateFile: "new-profile.html",
		Data: struct {
			Error string
		}{},
	}

	d.Render(w, r)
}

func (c *Context) profile(w web.ResponseWriter, r *web.Request) {
	id := r.URL.Query().Get("id")

	user := stats.GetUserBySlackID(id)
	for slack, su := range c.SlackUsers {
		if id == slack {
			user.Email = su.Profile.Email
		}
	}

	photos, err := recognition.GetUserPhotos(id)
	if err != nil {
		log.WithError(err).Error("unable to fin user photos")
	}

	d := tpl.Data{
		TemplateFile: "profile.html",
		Data: struct {
			Error  string
			User   *stats.User
			Photos []recognition.Photo
			Events []stats.Event
		}{
			User:   user,
			Photos: photos,
			Events: stats.GetUserEvents(id),
		},
	}

	d.Render(w, r)
}

func (c *Context) trainFile(w web.ResponseWriter, r *web.Request) {
	r.ParseMultipartForm(32 << 20)

	f, _, err := r.FormFile("file")
	slack := r.FormValue("slack")

	if err != nil {
		log.WithError(err).Error("unable to upload file")
	} else {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		if len(slack) > 0 {
			buf, readErr := ioutil.ReadAll(f)
			if readErr != nil {
				log.WithError(readErr).Error("unable to read file")
			} else {
				newFile, createErr := os.OpenFile(fmt.Sprintf("./faces-db/%s-%s.png", slack, id), os.O_WRONLY|os.O_CREATE, 0777)
				if createErr != nil {
					log.WithError(createErr).Error("unable to create file")
				} else {
					io.Copy(newFile, bytes.NewReader(buf))
					newFile.Close()
				}

				c.trainFacebox(bytes.NewReader(buf), id, slack)
			}

			f.Close()
		}
	}

	http.Redirect(w, r.Request, "/profile?id="+slack, http.StatusFound)
}

func (c *Context) trainFacebox(f io.Reader, id string, slack string) {
	trainErr := c.Fbox.Teach(f, id, slack)
	if trainErr != nil {
		log.WithError(trainErr).Error("unable to train")
	} else {
		log.WithField("slack", slack).Error("trained")
	}
}

func (c *Context) deletePhoto(w web.ResponseWriter, r *web.Request) {
	id := r.URL.Query().Get("id")
	slack := r.URL.Query().Get("slack")

	if len(id) > 0 {
		if err := c.Fbox.Remove(id); err != nil {
			log.WithError(err).WithField("id", id).Error("unable to remove face from facebox")
		} else {
			os.Remove(fmt.Sprintf("./faces-db/%s-%s.png", slack, id))
		}
	}

	http.Redirect(w, r.Request, "/profile?id="+slack, http.StatusFound)
}
