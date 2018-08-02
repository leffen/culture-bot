package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

// SpeechToText using Google Speech API
func SpeechToText(timeLimitSecs int) (string, error) {
	if mkdirErr := createFolderIfNotExists("record"); mkdirErr != nil {
		return "", mkdirErr
	}

	filename := fmt.Sprintf("record/%d.wav", time.Now().UnixNano())
	defer os.Remove(filename)

	if err := Record(filename, timeLimitSecs); err != nil {
		return "", err
	}

	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx, option.WithServiceAccountFile("speech.json"))
	if err != nil {
		return "", err
	}

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
			MaxAlternatives: 1,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return "", err
	}

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			return alt.Transcript, err
		}
	}

	return "", errors.New("no results found")
}
