package pin

import (
	"net/http"
	"net/url"
	"strings"
)

type PostsService struct {
	client *Client
}

type Post struct {
	Title       string
	Description string
	Hash        string
	URL         string
	Tags        []string
	ToRead      bool
}

func newPostFromPostResp(presp *postResp) *Post {
	var toRead bool
	if presp.ToRead == "yes" {
		toRead = true
	}

	return &Post{
		Title:       presp.Title,
		Description: presp.Description,
		Hash:        presp.Hash,
		URL:         presp.URL,
		Tags:        strings.Split(presp.Tag, " "),
		ToRead:      toRead,
	}
}

type postResp struct {
	Title       string `xml:"description,attr"`
	Description string `xml:"extended,attr"`
	Hash        string `xml:"hash,attr"`
	URL         string `xml:"href,attr"`
	Tag         string `xml:"tag,attr"`
	ToRead      string `xml:"toread,attr"`
}

// Recent fetches the most recent Posts for the authenticated account, filtered
// by tag. Optional filtering params can be provided in p.
//
// Valid params to pass are:
//
// * tag - up to 3 tags to filter by
// * count - number of results to return, default is 15, max is 100
//
// https://pinboard.in/api/#posts_recent
func (s *PostsService) Recent(p *url.Values) ([]*Post, *http.Response, error) {
	req, err := s.client.NewRequest("posts/recent", p)
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
