package pin

import (
	"net/http"
	"net/url"
)

// TagsService is the service for accessing Tag-related calls from the
// Pinboard API.
type TagsService struct {
	client *Client
}

type Tag struct {
	Count int    `xml:"count,attr"`
	Name  string `xml:"tag,attr"`
}

// Returns a full list of the user's tags along with the number of times they were used.
//
// https://pinboard.in/api#tags_get
func (s *TagsService) Get() ([]*Tag, *http.Response, error) {
	req, err := s.client.NewRequest("tags/get", nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Tags []*Tag `xml:"tag"`
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Tags, resp, nil
}

// Delete an existing tag.
//
// https://pinboard.in/api#tags_delete
func (s *TagsService) Delete(tag string) (*http.Response, error) {
	params := &url.Values{
		"tag": {tag},
	}
	req, err := s.client.NewRequest("tags/delete", params)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Rename an tag, or fold it in to an existing tag
//
// https://pinboard.in/api#tags_rename
func (s *TagsService) Rename(newTag, oldTag string) (*http.Response, error) {
	params := &url.Values{
		"old": {oldTag},
		"new": {newTag},
	}
	req, err := s.client.NewRequest("tags/rename", params)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
