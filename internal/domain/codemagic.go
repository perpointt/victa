package domain

import "time"

// CodemagicApplication описывает часть "application" ответа.
type CodemagicApplication struct {
	ID          string `json:"_id"`
	AppName     string `json:"appName"`
	ProjectType string `json:"projectType"`
	IconURL     string `json:"iconUrl"`
	LastBuildID string `json:"lastBuildId"`
}

type CodemagicBuild struct {
	ID         string    `json:"_id"`
	Status     string    `json:"status"`
	Version    string    `json:"version"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
	Commit     struct {
		AuthorName    string `json:"authorName"`
		CommitMessage string `json:"commitMessage"`
		Branch        string `json:"branch"`
	} `json:"commit"`
	Config struct {
		Name          string `json:"name"`
		BuildSettings struct {
			FlutterVersion string   `json:"flutterVersion"`
			Platforms      []string `json:"platforms"`
		} `json:"buildSettings"`
	} `json:"config"`
	BuildActions []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"buildActions"`
	Message   string `json:"message"`
	Artefacts []struct {
		Type      string `json:"type"`
		Path      string `json:"path"`
		PublicURL string `json:"public_url"`
	} `json:"artefacts"`
}

// CodemagicBuildResponse объединяет application + build
type CodemagicBuildResponse struct {
	Application CodemagicApplication `json:"application"`
	Build       CodemagicBuild       `json:"build"`
}
