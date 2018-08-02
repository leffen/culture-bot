package main

import "fmt"

// Step represents flow step
type Step struct {
	ID           string
	ValidateFunc func(s *Session) bool
	ExecuteFunc  func(s *Session) error
}

// Session represents user session data
type Session struct {
	End                      bool
	Greeted                  bool
	UserRepliedToBotGreeting bool
	IsCorrectPerson          bool
	IsWrongPerson            bool
	BotAskedCalendarQuestion bool
	UserWantsCalendar        bool
	UserDoesntWantCalendar   bool
	Face                     *RecognizeResponse
}

var flowSteps = []Step{
	Step{
		ID: "greeting",
		ValidateFunc: func(s *Session) bool {
			return !s.Greeted
		},
		ExecuteFunc: func(s *Session) error {
			speech := Speech{Folder: "audio", Language: "en"}
			if err := speech.Speak(s.Face.Message); err != nil {
				return err
			}

			s.Greeted = true
			return nil
		},
	},
	Step{
		ID: "get_user_reply_to_greeting",
		ValidateFunc: func(s *Session) bool {
			return s.Greeted && !s.UserRepliedToBotGreeting
		},
		ExecuteFunc: func(s *Session) error {
			text, sttErr := SpeechToText(3)
			if sttErr != nil {
				return sttErr
			}

			s.UserRepliedToBotGreeting = true
			s.IsCorrectPerson = isYes(text)
			if !s.IsCorrectPerson {
				s.IsWrongPerson = isNo(text)
			}
			return nil
		},
	},
	Step{
		ID: "greeting_yes",
		ValidateFunc: func(s *Session) bool {
			return s.Greeted && s.UserRepliedToBotGreeting && s.IsCorrectPerson && !s.BotAskedCalendarQuestion
		},
		ExecuteFunc: func(s *Session) error {
			speech := Speech{Folder: "audio", Language: "en"}
			if err := speech.Speak("Awesome! Shall I tell you your next meeting?"); err != nil {
				return err
			}

			s.BotAskedCalendarQuestion = true
			return nil
		},
	},
	Step{
		ID: "greeting_no",
		ValidateFunc: func(s *Session) bool {
			return s.Greeted && s.UserRepliedToBotGreeting && s.IsWrongPerson && !s.BotAskedCalendarQuestion
		},
		ExecuteFunc: func(s *Session) error {
			s.End = true

			speech := Speech{Folder: "audio", Language: "en"}
			return speech.Speak("Oops! Seems like my creators need to train me some more. See you!")
		},
	},
	Step{
		ID: "get_user_reply_to_calendar",
		ValidateFunc: func(s *Session) bool {
			return s.Greeted && s.UserRepliedToBotGreeting && s.BotAskedCalendarQuestion
		},
		ExecuteFunc: func(s *Session) error {
			text, sttErr := SpeechToText(3)
			if sttErr != nil {
				return sttErr
			}

			s.UserWantsCalendar = isYes(text)
			if !s.UserWantsCalendar {
				s.UserDoesntWantCalendar = isNo(text)
			}
			return nil
		},
	},
	Step{
		ID: "calendar_yes",
		ValidateFunc: func(s *Session) bool {
			return s.Greeted && s.UserWantsCalendar
		},
		ExecuteFunc: func(s *Session) error {
			s.End = true

			meeting, err := GetNextMeeting(*apiAddr, s.Face.User.Slack)
			if err != nil || meeting == nil {
				return err
			}

			speech := Speech{Folder: "audio", Language: "en"}
			return speech.Speak(meeting.Meeting)
		},
	},
	Step{
		ID: "calendar_no",
		ValidateFunc: func(s *Session) bool {
			return s.Greeted && s.UserDoesntWantCalendar
		},
		ExecuteFunc: func(s *Session) error {
			s.End = true

			speech := Speech{Folder: "audio", Language: "en"}
			return speech.Speak(fmt.Sprintf("No problem, I am always here to help you. Have a great %s!", s.Face.DayTime))
		},
	},
}
