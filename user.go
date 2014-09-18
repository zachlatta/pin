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

// SecretRSSKey returns the authenticated user's secret RSS key for viewing
// private RSS feeds.
func (s *UserService) SecretRSSKey() (string, *http.Response, error) {
	var result struct {
		XMLName xml.Name `xml:"result"`
		Secret  string   `xml:",chardata"`
	}

	req, err := s.client.NewRequest("user/secret", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return "", nil, err
	}

	return result.Secret, resp, nil
}
