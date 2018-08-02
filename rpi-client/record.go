package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Record from mic to a file using SoX
func Record(fileName string, timeLimitSecs int) error {
	log.WithField("file", fileName).Info("started audio recording")

	start := time.Now()

	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "trim", "0", strconv.Itoa(timeLimitSecs), "silence", "1", "0.1", "1%", "1", "0.5", "1%")

	env := os.Environ()
	env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		log.WithField("file", fileName).WithField("duration", time.Since(start).String()).Info("audio recorded")
	}

	return err
}
