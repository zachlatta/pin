package pin

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// PostsService is the service for accessing Post-related calls from the
// Pinboard API.
type PostsService struct {
	client *Client
}

const (
	timeLayoutFull  = "2006-01-02T15:04:05Z"
	timeLayoutShort = "2006-01-02"
)

// Post represents a post stored in Pinboard. Fields are transformed from the
// actual response to be a bit more sane. For example, description from the
// response is renamed to Title and the extended field is renamed to
// Description.
type Post struct {
	Title       string
	Description string
	Hash        string
	URL         string
	Tags        []string
	ToRead      bool
	Time        *time.Time
}

func newPostFromPostResp(presp *postResp) *Post {
	var toRead bool
	if presp.ToRead == "yes" {
		toRead = true
	}

	dt, _ := time.Parse(timeLayoutFull, presp.Time)

	return &Post{
		Title:       presp.Title,
		Description: presp.Description,
		Hash:        presp.Hash,
		URL:         presp.URL,
		Tags:        strings.Split(presp.Tag, " "),
		ToRead:      toRead,
		Time:        &dt,
	}
}

type postResp struct {
	Title       string `xml:"description,attr"`
	Description string `xml:"extended,attr"`
	Hash        string `xml:"hash,attr"`
	URL         string `xml:"href,attr"`
	Tag         string `xml:"tag,attr"`
	ToRead      string `xml:"toread,attr"`
	Time        string `xml:"time,attr"`
}

//
type Date struct {
	Count int
	Date  *time.Time
}

func newDateFromPostResp(dresp *dateResp) (*Date, error) {
	dt, err := time.Parse(timeLayoutShort, dresp.Date)
	if err != nil {
		return nil, err
	}
	c, err := strconv.Atoi(dresp.Count)
	if err != nil {
		return nil, err
	}

	return &Date{
		Count: c,
		Date:  &dt,
	}, nil
}

type dateResp struct {
	Count string `xml:"count,attr"`
	Date  string `xml:"date,attr"`
}

// Add creates a new Post for the authenticated account. urlStr and title are
// required.
//
// https://pinboard.in/api/#posts_add
func (s *PostsService) Add(urlStr, title, description string, tags []string,
	creationTime *time.Time, replace, shared,
	toread bool) (*http.Response, error) {
	var strTime string
	if creationTime != nil {
		strTime = creationTime.Format(timeLayoutFull)
	}

	params := &url.Values{
		"url":         {urlStr},
		"description": {title},
		"extended":    {description},
		"tags":        tags,
		"dt":          {strTime},
		"replace":     {fmt.Sprintf("%t", replace)},
		"shared":      {fmt.Sprintf("%t", shared)},
		"toread":      {fmt.Sprintf("%t", toread)},
	}

	req, err := s.client.NewRequest("posts/add", params)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete deletes the specified Post from the authenticated account where
// urlStr is the URL of the Post to delete.
//
// https://pinboard.in/api/#posts_delete
func (s *PostsService) Delete(urlStr string) (*http.Response, error) {
	params := &url.Values{"url": {urlStr}}
	req, err := s.client.NewRequest("posts/delete", params)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Get returns one or more posts on a single day matching the arguments.
// If no date or url is given, date of most recent bookmark will be used.
//
// https://pinboard.in/api#posts_get
func (s *PostsService) Get(tags []string, creationTime *time.Time, urlStr string) ([]*Post, *http.Response, error) {

	params := &url.Values{}

	if creationTime != nil {
		params.Add("dt", creationTime.Format(timeLayoutFull))
	}

	if tags != nil && len(tags) > 3 {
		return nil, nil, errors.New("too many tags (max is 3)")
	} else if tags != nil && len(tags) > 0 {
		params.Add("tags", strings.Join(tags, " "))
	}

	if len(urlStr) > 0 {
		params.Add("url", urlStr)
	}

	req, err := s.client.NewRequest("posts/get", params)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Posts []*postResp `xml:"post"`
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	posts := make([]*Post, len(result.Posts))
	for i, v := range result.Posts {
		posts[i] = newPostFromPostResp(v)
	}

	return posts, resp, nil
}

// Returns the most recent time a bookmark was added, updated or deleted.
// Use this before calling posts/all to see if the data has changed since the last fetch.
//
// https://pinboard.in/api#posts_update
func (s *PostsService) LastTimeUpdated() (*time.Time, *http.Response, error) {
	req, err := s.client.NewRequest("posts/update", &url.Values{})
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Time string `xml:"time,attr"`
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	updated, err := time.Parse(timeLayoutFull, result.Time)
	if err != nil {
		return nil, resp, err
	}

	return &updated, resp, nil
}

// Returns a list of dates with the number of posts at each date.
//
// https://pinboard.in/api#posts_dates
func (s *PostsService) Dates(tags []string) ([]*Date, *http.Response, error) {
	params := &url.Values{}

	if tags != nil && len(tags) > 3 {
		return nil, nil, errors.New("too many tags (max is 3)")
	} else if tags != nil && len(tags) > 0 {
		params.Add("tags", strings.Join(tags, " "))
	}

	req, err := s.client.NewRequest("posts/dates", params)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Dates []*dateResp `xml:"date"`
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	dates := make([]*Date, len(result.Dates))
	for i, v := range result.Dates {
		d, err := newDateFromPostResp(v)
		if err != nil {
			return nil, nil, err
		}
		dates[i] = d
	}

	return dates, resp, nil
}

// Recent fetches the most recent Posts for the authenticated account, filtered
// by tag. Up to 3 tags can be specified to filter by. The max count is 100. If
// a negative count is passed, then the default number of posts (15) is
// returned.
//
// https://pinboard.in/api/#posts_recent
func (s *PostsService) Recent(tags []string, count int) ([]*Post,
	*http.Response, error) {
	if tags != nil && len(tags) > 3 {
		return nil, nil, errors.New("too many tags (max is 3)")
	}
	if count > 100 {
		return nil, nil, errors.New("count must be below 100")
	}
	if count < 0 {
		count = 15
	}

	req, err := s.client.NewRequest("posts/recent", &url.Values{
		"tag":   tags,
		"count": {strconv.Itoa(count)},
	})
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Posts []*postResp `xml:"post"`
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	posts := make([]*Post, len(result.Posts))
	for i, v := range result.Posts {
		posts[i] = newPostFromPostResp(v)
	}

	return posts, resp, nil
}

// All fetches all bookmarks in the user's account.
//
// https://pinboard.in/api#posts_all
func (s *PostsService) All(tags []string, start int, results int, fromdt, todt *time.Time) ([]*Post,
	*http.Response, error) {

	params := &url.Values{}

	if tags != nil && len(tags) > 3 {
		return nil, nil, errors.New("too many tags (max is 3)")
	} else if tags != nil && len(tags) > 0 {
		params.Add("tags", strings.Join(tags, " "))
	}

	if start > 0 {
		params.Add("start", strconv.Itoa(start))
	}

	if results > 0 {
		params.Add("results", strconv.Itoa(results))
	}

	if fromdt != nil {
		params.Add("fromdt", fromdt.Format(timeLayoutFull))
	}

	if todt != nil {
		params.Add("todt", todt.Format(timeLayoutFull))
	}

	req, err := s.client.NewRequest("posts/all", params)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Posts []*postResp `xml:"post"`
	}

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	posts := make([]*Post, len(result.Posts))
	for i, v := range result.Posts {
		posts[i] = newPostFromPostResp(v)
	}

	return posts, resp, nil

}

// TODO
//
// https://pinboard.in/api#posts_suggest
func (s *PostsService) Suggest() {
}
