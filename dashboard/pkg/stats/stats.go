package stats

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// User type
type User struct {
	First              string
	Last               string
	Slack              string
	Pronounce          string
	Email              string
	Photos             []string
	AverageRecognition int
}

// Event type
type Event struct {
	Slack  string
	Key    string
	Value  string
	At     string
	AtTime time.Time
}

var (
	db *sql.DB
)

// Init stats DB
func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./stats.db")
	if err != nil {
		log.WithError(err).Error("can't init db")
	} else {
		// Migration
		//db.Exec("create table logs (slack text not null, event_key text not null, value text not null, at datetime default CURRENT_TIMESTAMP)")
	}

	log.Info("sqlite ready")
}

// GetUserBySlackID returns user obj
func GetUserBySlackID(id string) *User {
	stmt, err := db.Prepare("select first, last, slack, pronounce from users where slack = ?")
	if err != nil {
		log.WithError(err).Error("can't select user")
		return nil
	}
	defer stmt.Close()

	var (
		first     string
		last      string
		slack     string
		pronounce string
	)
	err = stmt.QueryRow(id).Scan(&first, &last, &slack, &pronounce)
	if err != nil {
		log.WithError(err).Error("can't query user")
		return nil
	}

	return &User{
		First:     first,
		Last:      last,
		Slack:     slack,
		Pronounce: pronounce,
	}
}

// UpdateUser func
func UpdateUser(id string, u User) {
	tx, _ := db.Begin()

	stmt, err := tx.Prepare("update users set first = ?, last = ?, slack = ?, pronounce = ? where slack = ?")
	if err != nil {
		log.WithError(err).Error("can't create stmt")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.First, u.Last, u.Slack, u.Pronounce, id)
	if err != nil {
		log.WithError(err).Error("can't exec stmt")
	}

	tx.Commit()

	log.WithField("slack", id).Info("updated user")
}

// AddUser func
func AddUser(u User) {
	tx, _ := db.Begin()

	stmt, err := tx.Prepare("insert into users (first, last, slack, pronounce) values (?, ?, ?, ?)")
	if err != nil {
		log.WithError(err).Error("can't create stmt")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.First, u.Last, u.Slack, u.Pronounce)
	if err != nil {
		log.WithError(err).Error("can't exec stmt")
	}

	tx.Commit()

	log.WithField("slack", u.Slack).Info("added user")
}

// DeleteUser func
func DeleteUser(id string) {
	tx, _ := db.Begin()

	stmt, err := tx.Prepare("delete from users where slack = ?")
	if err != nil {
		log.WithError(err).Error("can't create stmt")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.WithError(err).Error("can't exec stmt")
	}

	tx.Commit()

	log.WithField("slack", id).Info("deleted user")
}

// GetUsers func
func GetUsers() []*User {
	rows, err := db.Query("select first, last, slack, pronounce from users")
	if err != nil {
		log.WithError(err).Error("can't get users")
		return nil
	}
	defer rows.Close()

	res := []*User{}
	for rows.Next() {
		var (
			first     string
			last      string
			slack     string
			pronounce string
		)
		err = rows.Scan(&first, &last, &slack, &pronounce)
		if err != nil {
			log.WithError(err).Error("scan error")
			continue
		}
		res = append(res, &User{
			First:     first,
			Last:      last,
			Slack:     slack,
			Pronounce: pronounce,
		})
	}

	return res
}

// LogEvent func
func LogEvent(e Event) {
	tx, _ := db.Begin()

	stmt, err := tx.Prepare("insert into logs (slack, event_key, value) values (?, ?, ?)")
	if err != nil {
		log.WithError(err).Error("can't create stmt")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Slack, e.Key, e.Value)
	if err != nil {
		log.WithError(err).Error("can't exec stmt")
	}

	tx.Commit()

	log.WithField("event", e).Info("saved log event")
}

// GetUserEvents func
func GetUserEvents(slack string) []Event {
	stmt, err := db.Prepare("select event_key, value, at from logs where slack = ? order by at desc")
	if err != nil {
		log.WithError(err).Error("can't prepare stmt")
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(slack)
	if err != nil {
		log.WithError(err).Error("can't get users")
		return nil
	}
	defer rows.Close()

	res := []Event{}
	for rows.Next() {
		var (
			key   string
			value string
			at    string
		)
		err = rows.Scan(&key, &value, &at)
		if err != nil {
			log.WithError(err).Error("scan error")
			continue
		}

		atTime, timeErr := time.Parse(time.RFC3339, at)
		if timeErr != nil {
			log.WithError(timeErr).Error("unable to parse time")
		}

		res = append(res, Event{
			Slack:  slack,
			Key:    key,
			Value:  value,
			At:     at,
			AtTime: atTime,
		})
	}

	return res
}

// GetAverageRecognitionValue func
func GetAverageRecognitionValue(slack string) int {
	count := 0
	sum := 0

	events := GetUserEvents(slack)
	for _, e := range events {
		if e.Key == "RECOGNITION" {
			valInt, _ := strconv.Atoi(e.Value)
			if valInt > 0 {
				count++
				sum += valInt
			}
		}
	}

	if count == 0 {
		return 0
	}

	return int(sum / count)
}
