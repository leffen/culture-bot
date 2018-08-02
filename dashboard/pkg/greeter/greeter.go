package greeter

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/plutov/culture-bot/dashboard/pkg/stats"
)

// greeting
type greeting struct {
	Time time.Time
}

// Face : recognized face
type Face struct {
	Slack string
}

var (
	greetings = make(map[string]greeting)
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// IsGreeted : bool
func IsGreeted(f Face) bool {
	g, ok := greetings[f.Slack]
	now := time.Now()
	return ok && now.Before(g.Time.Add(time.Hour*12))
}

// MarkAsGreeted : say to user
func MarkAsGreeted(f Face, tz string) (string, string, *stats.User) {
	greetings[f.Slack] = greeting{
		Time: time.Now(),
	}

	log.WithField("name", f.Slack).Info("greeting")
	u := stats.GetUserBySlackID(f.Slack)
	if u == nil {
		return "", "", nil
	}

	dt := GetDayTimeInfo(tz)

	return fmt.Sprintf("Good %s. Is that you %s?", dt, u.Pronounce), dt, u
}

// ClearGreetings func is executed when client is restarted
func ClearGreetings() {
	greetings = make(map[string]greeting)
}

// GetDayTimeInfo returns morning|afternoon|evening
func GetDayTimeInfo(tz string) string {
	now := time.Now()
	if len(tz) > 0 {
		loc, locErr := time.LoadLocation(tz)
		if locErr != nil {
			log.WithError(locErr).Error("unable to load location")
		} else {
			now = time.Now().In(loc)
		}
	}

	dayTimeStr := "evening"
	if now.Hour() >= 4 && now.Hour() < 12 {
		dayTimeStr = "morning"
	} else if now.Hour() >= 12 && now.Hour() < 18 {
		dayTimeStr = "afternoon"
	}

	return dayTimeStr
}
