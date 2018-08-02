package main

import (
	log "github.com/sirupsen/logrus"
)

// StartSession starts conversation with user
func StartSession(face *RecognizeResponse) {
	log.WithField("user", face.User.Slack).Info("session started")
	defer log.WithField("user", face.User.Slack).Info("session finished")

	session := &Session{
		Face: face,
	}

	for {
		var stepFound bool
		for _, step := range flowSteps {
			if step.ValidateFunc(session) {
				stepFound = true

				log.WithFields(log.Fields{"step": step.ID, "slack": session.Face.User.Slack}).Info("executing step")
				if err := step.ExecuteFunc(session); err != nil {
					log.WithError(err).WithFields(log.Fields{"step": step.ID, "slack": session.Face.User.Slack}).Error("step execution error")
					return
				}
			}
		}

		// End Session
		if !stepFound || session.End {
			return
		}
	}
}
