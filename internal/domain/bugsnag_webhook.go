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
	Error struct {
		Message       string     `json:"message"`
		URL           string     `json:"url"`
		Status        string     `json:"status"`
		Unhandled     bool       `json:"unhandled"`
		Occurrences   int64      `json:"occurrences"`
		FirstReceived *time.Time `json:"firstReceived"`
		ReceivedAt    *time.Time `json:"receivedAt"`
		UserID        string     `json:"userId"`
		App           struct {
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
