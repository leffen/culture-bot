package greeter

import (
	"testing"
	"time"
)

func TestGetDayTimeInfo(t *testing.T) {
	now := time.Now()

	hcmc := GetDayTimeInfo("Asia/Ho_Chi_Minh")
	bangkok := GetDayTimeInfo("Asia/Bangkok")
	local := GetDayTimeInfo("")
	la := GetDayTimeInfo("America/Los_Angeles")

	if now.Hour() >= 4 && now.Hour() < 12 && bangkok != "morning" {
		t.Errorf("4-12: morning")
	}

	if now.Hour() >= 12 && now.Hour() < 18 && bangkok != "afternoon" {
		t.Errorf("12-18: afternoon")
	}

	if (now.Hour() >= 18 || now.Hour() < 4) && bangkok != "evening" {
		t.Errorf("18-4: evening")
	}

	if hcmc != local {
		t.Errorf("hcmc is local")
	}

	if hcmc != bangkok {
		t.Errorf("hcmc is bangkok")
	}

	if la == local {
		t.Errorf("la is not local")
	}
}
