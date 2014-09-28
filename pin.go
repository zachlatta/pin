package pin

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://api.pinboard.in/v1/"
	userAgent      = "pin/" + libraryVersion
)

type AuthToken struct {
	Username string
	Token    string
}

func (t *AuthToken) String() string {
	return fmt.Sprintf("%s:%s", t.Username, t.Token)
}

type Client struct {
	client    *http.Client
	authToken *AuthToken
	BaseURL   *url.URL
	UserAgent string

	Posts *PostsService
	Tags  *TagsService
	User  *UserService
	Notes *NotesService
}

// NewClient returns a new Pinboard API client. If a nil httpClient client is
// provided, http.DefaultClient will be used. If you want to do authenticated
// requests using HTTP Auth, you need to pass in an authenticated client - this
// library does not handle authentication.
func NewClient(httpClient *http.Client, authToken *AuthToken) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		authToken: authToken,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.Posts = &PostsService{c}
	c.Tags = &TagsService{c}
	c.User = &UserService{c}
	c.Notes = &NotesService{c}
	return c
}

// NewRequest constructs a new request to the Pinboard API. A relative URL can
// be provided in urlStr, in which case it's resolved to the Client's BaseURL.
// Relative URLs should always be specified without a preceding slash. If the
// Client has an AuthToken set, then it is added to the urlParams.
func (c *Client) NewRequest(urlStr string,
	urlParams *url.Values) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if c.authToken != nil {
		tok := c.authToken.String()
		if urlParams != nil {
			urlParams.Add("auth_token", tok)
		} else {
			urlParams = &url.Values{"auth_token": {tok}}
		}
	}

	u := c.BaseURL.ResolveReference(rel)
	u.RawQuery = urlParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// XML decoded and stored in the value pointed to by v, or returned as an error
// if an API error has occured. If v implements the io.Writer interface, the
// raw response will be written to v, without attempting to first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = xml.NewDecoder(resp.Body).Decode(v)
		}
	}
	return resp, err
}
