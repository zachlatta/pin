package pin

import (
	"encoding/xml"
	"net/http"
)

// UserService provides methods for accessing user actions through the Pinboard
// API.
type UserService struct {
	client *Client
}

type xmlResult struct {
	XMLName xml.Name `xml:"result"`
	Body    string   `xml:",chardata"`
}

// SecretRSSKey returns the authenticated user's secret RSS key for viewing
// private RSS feeds.
func (s *UserService) SecretRSSKey() (string, *http.Response, error) {
	result := xmlResult{}
	req, err := s.client.NewRequest("user/secret", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return "", nil, err
	}

	return result.Body, resp, nil
}

// APIToken returns the authenticated user's API token.
func (s *UserService) APIToken() (string, *http.Response, error) {
	result := xmlResult{}
	req, err := s.client.NewRequest("user/api_token", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return "", nil, err
	}

	return result.Body, resp, nil
}
