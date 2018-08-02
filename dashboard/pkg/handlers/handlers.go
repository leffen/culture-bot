package handlers

import (
	"net/http"
	"os"
	"path"

	"github.com/gocraft/web"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	sl "github.com/plutov/culture-bot/dashboard/pkg/slack"
)

// Context struct
type Context struct {
	Fbox       *facebox.Client
	Slack      *slack.Client
	SlackUsers map[string]slack.User
}

// ListenAndServe func
func ListenAndServe() error {
	ctx := new(Context)
	ctx.Fbox = facebox.New("http://facebox:8080")
	ctx.Slack = sl.New()
	ctx.SlackUsers = sl.GetAllUsers(ctx.Slack)

	currentRoot, _ := os.Getwd()

	r := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.StaticMiddleware(path.Join(currentRoot, "faces-db"))).
		Get("/", ctx.home).
		Get("/people", ctx.people).
		Get("/settings", ctx.settings).
		Get("/profile", ctx.profile).
		Get("/profile/new", ctx.newProfile).
		Get("/analytics/delete/user", ctx.deleteUser).
		Get("/analytics/delete/photo", ctx.deletePhoto).
		Post("/analytics/update/user", ctx.updateUser).
		Post("/analytics/add/user", ctx.addUser).
		Get("/client", ctx.client).
		Post("/train/file", ctx.trainFile).
		Post("/recognize", ctx.recognize).
		Post("/rpi/recognize", ctx.rpiRecognize).
		Post("/rpi/clear-greetings", ctx.rpiClearGreetings).
		Post("/calendar", ctx.calendar).
		Get("/calendar/token", ctx.calendarToken).
		Get("/calendar/redirect", ctx.calendarRedirect)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", r)

	log.Info("web service started")

	return http.ListenAndServe("0.0.0.0:8081", serveMux)
}
