package domain

import "time"

type BugsnagWebhook struct {
	Project struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"project"`
	Trigger struct {
		Type string `json:"type"`
	} `json:"trigger"`
	UserID        string     `json:"user_id"`
	FirstReceived *time.Time `json:"firstReceived"`
	ReceivedAt    *time.Time `json:"receivedAt"`
	Severity      string     `json:"severity"`
	Occurrences   int64      `json:"occurrences"`
	Error         struct {
		Message string `json:"message"`
		Context string `json:"context"`
		URL     string `json:"url"`
		App     struct {
			ID          string `json:"id"`
			Version     string `json:"version"`
			VersionCode string `json:"versionCode"`
			Type        string `json:"type"`
		} `json:"app"`
		Device struct {
			ID           string    `json:"id"`
			Manufacturer string    `json:"manufacturer"`
			Model        string    `json:"model"`
			OSName       string    `json:"osName"`
			OSVersion    string    `json:"osVersion"`
			FreeMemory   int64     `json:"freeMemory"`
			TotalMemory  int64     `json:"totalMemory"`
			FreeDisk     int64     `json:"freeDisk"`
			JailBroken   bool      `json:"jailbroken"`
			Orientation  string    `json:"orientation"`
			Locale       string    `json:"locale"`
			Charging     bool      `json:"charging"`
			BatteryLevel float64   `json:"batteryLevel"`
			Time         time.Time `json:"time"`
		} `json:"device"`
		Exceptions []struct {
			Message    string `json:"message"`
			StackTrace []struct {
				File       string `json:"file"`
				LineNumber string `json:"lineNumber"`
				Method     string `json:"method"`
			} `json:"stacktrace"`
		} `json:"exceptions"`
	} `json:"error"`
}
