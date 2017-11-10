package spotilocal

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

const (
	DESKTOP_CLIENT_URL = "https://localhost.spotilocal.com:%d"
)

var (
	ErrNotRunning   = fmt.Errorf("Spotify Desktop Client is not running")
	ErrDisconnected = fmt.Errorf("Disconnected from Spotify Desktop Client")
)

// Spotify is a struct to perform actions on Spotify Desktop Client
type Spotify struct {
	token string
	csrf  string
	port  int
}

// New returns Spotify instance
func New() (*Spotify, error) {
	token, err := fetchToken()

	if err != nil {
		return nil, err
	}

	return NewWithToken(token)
}

// NewWithToken returns Spotify instance with token
func NewWithToken(token string) (*Spotify, error) {
	return &Spotify{
		token: token,
	}, nil
}

// Connect finds Spotify Desktop Client and fetches CSRF token
func (s *Spotify) Connect() error {
	for port := 4370; port <= 4400; port++ {
		csrfToken, err := s.csrfToken(port)
		if err != nil {
			if _, ok := err.(*net.OpError); !ok {
				return err
			}

			// It must be some connection error (*net.OpError). Lets check next port.
			continue
		}

		// Since theres no error, that means it might be the port we are looking for
		s.port = port
		s.csrf = csrfToken
		return nil
	}

	return ErrNotRunning
}

func (s *Spotify) requestWithPort(path string, query map[string]string, port int) ([]byte, error) {
	// Prepare request
	baseURL := fmt.Sprintf(DESKTOP_CLIENT_URL, port)
	completeURL := baseURL + path
	req, err := http.NewRequest("GET", completeURL, nil)

	// Add query string
	q := req.URL.Query()
	q.Add("oauth", s.token)
	if s.csrf != "" {
		q.Add("csrf", s.csrf)
	}

	// Extend query string
	if query != nil {
		for key, val := range query {
			q.Add(key, val)
		}
	}
	req.URL.RawQuery = q.Encode()

	// Set Origin header
	req.Header.Add("Origin", ORIGIN_URL)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Handle EOF error. Usually EOF occurs when spotify client is closed.
		urlErr, ok := err.(*url.Error)
		if ok && urlErr.Err == io.EOF {
			s.port = 0
			s.csrf = ""
			return nil, ErrDisconnected
		}
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *Spotify) request(path string, query map[string]string) ([]byte, error) {
	if s.port == 0 {
		return nil, ErrDisconnected
	}

	return s.requestWithPort(path, query, s.port)
}
