package spotilocal

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// APIs for spotilocal
const (
	CSRF_PATH    = "/simplecsrf/token.json"
	PLAY_PATH    = "/remote/play.json"
	PAUSE_PATH   = "/remote/pause.json"
	STATUS_PATH  = "/remote/status.json"
	VERSION_PATH = "/service/version.json"
)

// Event is used for getting status on particular events
type Event string

// Enums for Events
const (
	OnAP     Event = "ap"
	OnError  Event = "error"
	OnLogin  Event = "login"
	OnLogout Event = "logout"
	OnPause  Event = "pause"
	OnPlay   Event = "play"
)

var (
	ErrEmptyResponse = fmt.Errorf("Empty response")
)

// csrfToken returns CSRF Token from a given server (port)
func (s *Spotify) csrfToken(port int) (string, error) {
	buf, err := s.requestWithPort(CSRF_PATH, nil, port)

	if err != nil {
		urlErr, ok := err.(*url.Error)
		if !ok {
			return "", err
		}

		opErr, ok := urlErr.Err.(*net.OpError)
		if !ok {
			return "", urlErr
		}

		return "", opErr
	}

	body := struct {
		Token string `json:"token"`
	}{}
	if err := json.Unmarshal(buf, &body); err != nil {
		return "", err
	}

	return body.Token, nil
}

// Pause sets the pause status
func (s *Spotify) Pause(status bool) error {
	statusStr := "true"
	if !status {
		statusStr = "false"
	}

	query := map[string]string{
		"pause": statusStr,
	}

	if _, err := s.request(PAUSE_PATH, query); err != nil {
		return err
	}

	return nil
}

// Play starts playing requested track
func (s *Spotify) Play(uri string) error {
	return s.PlayWithContext(uri, uri)
}

// PlayWithContext starts playing requested track in context with an album or artist
func (s *Spotify) PlayWithContext(uri, context string) error {
	query := map[string]string{
		"uri":     uri,
		"context": context,
	}

	if _, err := s.request(PLAY_PATH, query); err != nil {
		return err
	}

	return nil
}

// Status gets the current status immediately
func (s *Spotify) Status() (Status, error) {
	return s.StatusOn()
}

// StatusOn gets the current status on a particular Event
// Providing an empty event array will return the current status
// Incase of empty response, it will throw ErrEmptyResponse error
//
// Spotilocal client ignores events when it is not yet properly
// spin up and causes request to respond immediately with empty
// status values.
func (s *Spotify) StatusOn(events ...Event) (Status, error) {
	eventsStr := []string{}
	for _, e := range events {
		eventsStr = append(eventsStr, string(e))
	}

	query := map[string]string{
		"returnon": strings.Join(eventsStr, ","),
	}

	buf, err := s.request(STATUS_PATH, query)
	if err != nil {
		return Status{}, err
	}

	body := Status{}
	if err := json.Unmarshal(buf, &body); err != nil {
		return Status{}, err
	}

	if body == (Status{}) {
		return Status{}, ErrEmptyResponse
	}

	return body, nil
}

// Version gets the version and client version
func (s *Spotify) Version() (int, string, error) {
	query := map[string]string{
		"service": "remote",
	}

	buf, err := s.request(VERSION_PATH, query)
	if err != nil {
		return 0, "", err
	}

	body := struct {
		Version       int    `json:"version"`
		ClientVersion string `json:"client_version"`
	}{}
	if err := json.Unmarshal(buf, &body); err != nil {
		return 0, "", err
	}

	return body.Version, body.ClientVersion, nil
}
