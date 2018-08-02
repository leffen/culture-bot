package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// User structure same as in API
type User struct {
	First     string
	Last      string
	Slack     string
	Pronounce string
}

// RecognizeResponse :api response
type RecognizeResponse struct {
	Message string
	DayTime string
	User    *User
}

// CalendarResponse : ajax response
type CalendarResponse struct {
	Email   string
	Meeting string
}

// Recognize faces from frame
func Recognize(apiAddr string, frame []byte) (*RecognizeResponse, error) {
	sEnc := base64.StdEncoding.EncodeToString(frame)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", apiAddr, "/rpi/recognize"), bytes.NewBufferString(sEnc))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	res := new(RecognizeResponse)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ClearGreetings when program is started
func ClearGreetings(apiAddr string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", apiAddr, "/rpi/clear-greetings"), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetNextMeeting gets next user meeting
func GetNextMeeting(apiAddr string, slack string) (*CalendarResponse, error) {
	form := url.Values{}
	form.Add("slack", slack)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", apiAddr, "/calendar"), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	res := new(CalendarResponse)
	err = json.NewDecoder(resp.Body).Decode(res)

	return res, err
}
