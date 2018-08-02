package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gocraft/web"
	log "github.com/sirupsen/logrus"
	"github.com/plutov/culture-bot/dashboard/pkg/greeter"
	"github.com/plutov/culture-bot/dashboard/pkg/recognition"
	"github.com/plutov/culture-bot/dashboard/pkg/stats"
)

// RecognizeResponse :api response
type RecognizeResponse struct {
	Message string
	DayTime string
	User    *stats.User
}

func (c *Context) recognize(w web.ResponseWriter, r *web.Request) {
	input := r.FormValue("img")

	// Normalize image from web canvas
	b64data := input[strings.IndexByte(input, ',')+1:]
	buf, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.WithError(err).Error("unable to decode base64")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := c.handleWebCamImage(buf)

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (c *Context) rpiRecognize(w web.ResponseWriter, r *web.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("unable to read body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.WithError(err).Error("unable to decode base64")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := c.handleWebCamImage(buf)

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (c *Context) rpiClearGreetings(w web.ResponseWriter, r *web.Request) {
	defer r.Body.Close()

	greeter.ClearGreetings()
}

func (c *Context) handleWebCamImage(buf []byte) *RecognizeResponse {
	faces, err := c.Fbox.Check(bytes.NewReader(buf))
	if err != nil {
		log.WithError(err).Error("unable to recognize face")
		return nil
	}

	sort.Sort(recognition.Faces(faces))

	var resp *RecognizeResponse

	for _, f := range faces {
		log.WithField("face", f).Info("found")

		if len(f.Name) > 0 && f.Confidence >= 0.7 {
			faceObj := greeter.Face{
				Slack: f.Name,
			}

			// Only greet first ungreeted person
			if resp == nil && !greeter.IsGreeted(faceObj) {
				su, ok := c.SlackUsers[f.Name]
				tz := ""
				if ok {
					tz = su.TZ
				}
				msg, dt, u := greeter.MarkAsGreeted(faceObj, tz)
				stats.LogEvent(stats.Event{
					Slack: f.Name,
					Key:   "GREETING",
				})

				resp = &RecognizeResponse{
					Message: msg,
					DayTime: dt,
					User:    u,
				}
			}
		}

		if len(f.Name) > 0 && f.Confidence > 0 {
			stats.LogEvent(stats.Event{
				Slack: f.Name,
				Key:   "RECOGNITION",
				Value: fmt.Sprintf("%d", int(f.Confidence*100)),
			})

			if saveErr := recognition.SaveFaceImage(f, buf, f.Name); saveErr != nil {
				log.WithError(saveErr).Error("unable to save face into a file")
			}
		}
	}

	return resp
}
