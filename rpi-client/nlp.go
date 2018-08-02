package main

import "strings"

func isYes(text string) bool {
	keywords := []string{"yes", "sure", "ok", "yea", "yep", "yea", "it is", "yeah", "yup", "ya", "correct", "right", "exactly", "definitely", "why not", "thats me"}
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}

	return false
}

func isNo(text string) bool {
	keywords := []string{"no", "nope", "oops", "not me", "nah", "no thanks", "not"}
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}

	return false
}
