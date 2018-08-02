package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocraft/web"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/plutov/culture-bot/dashboard/pkg/greeter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

// CalendarResponse : ajax response
type CalendarResponse struct {
	Email   string
	Meeting string
}

func (c *Context) calendar(w web.ResponseWriter, r *web.Request) {
	r.ParseForm()

	var email string
	slackName := r.FormValue("slack")

	su, ok := c.SlackUsers[slackName]
	if ok {
		email = su.Profile.Email
	}

	cr := CalendarResponse{
		Email: email,
	}

	var (
		err      error
		msgNoLoc string
	)

	if len(email) > 0 {
		cr.Meeting, msgNoLoc, err = getNextMeeting(email)
		if err == nil {
			log.WithField("meeting", cr.Meeting).Info("sent meeting to client")

			// Also send to Slack
			go func() {
				if _, _, slackErr := c.Slack.PostMessage(su.ID, msgNoLoc, slack.PostMessageParameters{
					AsUser: true,
				}); slackErr != nil {
					log.WithError(slackErr).Error("unable to post to slack")
				}
			}()
		}
	}

	if len(cr.Meeting) == 0 {
		log.WithError(err).WithField("slack", slackName).Error("unable to get next meeting")
		cr.Meeting = "Sorry, I couldn't find any meetings in your Calendar."
	}

	js, _ := json.Marshal(cr)
	w.Write(js)
}

func getNextMeeting(email string) (string, string, error) {
	client := newOAuthClient()
	if client == nil {
		return "", "", fmt.Errorf("client is nil")
	}

	service, err := calendar.New(client)
	if err != nil {
		return "", "", err
	}

	events, err := service.Events.List(email).ShowDeleted(false).
		SingleEvents(true).TimeMin(time.Now().Format(time.RFC3339)).MaxResults(5).OrderBy("startTime").Do()
	if err != nil {
		return "", "", err
	}

	if len(events.Items) < 1 {
		return "", "", fmt.Errorf("no events found")
	}

	msg := ""
	msgNoLoc := ""
	for _, i := range events.Items {
		if i.Start == nil {
			continue
		}

		t := i.Start

		loc, locErr := time.LoadLocation(t.TimeZone)
		if locErr != nil {
			log.WithError(locErr).Error("unable to load location")
			continue
		}

		nowTZ := time.Now().In(loc)

		startTime, parseErr := time.Parse(time.RFC3339, t.DateTime)
		if parseErr != nil {
			log.WithError(parseErr).WithField("date", t.DateTime).Error("unable to parse date")
			continue
		}

		if i.Kind == "calendar#event" && startTime.After(nowTZ) {
			msg = fmt.Sprintf("Your next meeting is \"%s\", starting at %s", i.Summary, startTime.Format(time.Kitchen))
			msgNoLoc = msg
			where := getLocationName(i)
			if len(where) > 0 {
				msg += fmt.Sprintf(", in %s", where)
			}
			msg += fmt.Sprintf(". I've sent it to you via Slack. Have a great %s!", greeter.GetDayTimeInfo(t.TimeZone))
			break
		}
	}

	if len(msg) == 0 {
		return "", "", fmt.Errorf("no events found")
	}

	return msg, msgNoLoc, nil
}

func getLocationName(e *calendar.Event) string {
	if len(e.Location) == 0 {
		return ""
	}

	// Use last room from list separated by comma
	parts := strings.Split(e.Location, ",")
	e.Location = strings.TrimSpace(parts[len(parts)-1])

	vnMapping := map[string]string{
		"Tạ Quang VN (8)": "Ta Qwan",
		"Truong VN (7)":   "Tchoeng",
		"Trịnh VN (6)":    "Chin",
	}

	replaceVN, ok := vnMapping[e.Location]
	if ok {
		e.Location = replaceVN
	}

	reg := regexp.MustCompile("[^a-zA-Z ]+")
	return strings.TrimSpace(reg.ReplaceAllString(e.Location, ""))
}

func (c *Context) calendarToken(w web.ResponseWriter, r *web.Request) {
	config := getConfig()

	randState := fmt.Sprintf("state-%d", time.Now().UnixNano())
	authURL := config.AuthCodeURL(randState)

	log.WithField("url", authURL).Info("redirecting to auth url")
	http.Redirect(w, r.Request, authURL, http.StatusFound)
}

func (c *Context) calendarRedirect(w web.ResponseWriter, r *web.Request) {
	config := getConfig()
	ctx := context.Background()

	if code := r.FormValue("code"); code != "" {
		log.WithField("code", code).Error("recieved code")

		token, err := config.Exchange(ctx, code)
		if err != nil {
			log.WithError(err).Error("token exchange error")
			http.Error(w, "", 500)
		} else {
			saveToken(token)
			http.Redirect(w, r.Request, os.Getenv("NGROK_ADDR"), http.StatusFound)
		}
		return
	}

	log.Error("no code")
	http.Error(w, "", 500)
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CALENDAR_CLEINT_ID"),
		ClientSecret: os.Getenv("CALENDAR_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("NGROK_ADDR") + "/calendar/redirect",
		Scopes:       []string{calendar.CalendarReadonlyScope},
	}
}

func newOAuthClient() *http.Client {
	config := getConfig()
	ctx := context.Background()

	token, err := tokenFromFile()
	if err != nil {
		log.WithError(err).Error("token not found")
		return nil
	}

	return config.Client(ctx, token)
}

func tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open("calendar_token.dat")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(token *oauth2.Token) {
	f, err := os.Create("calendar_token.dat")
	if err != nil {
		log.WithError(err).Error("failed to cache oauth token")
		return
	}

	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}
