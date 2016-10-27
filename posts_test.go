package pin

import (
	"io/ioutil"
	"strings"
	"testing"

	"time"

	"github.com/jarcoal/httpmock"
)

var (
	token  = AuthToken{Username: "user", Token: "token"}
	client = NewClient(nil, &token)
)

func readFixture(filename string) string {
	data, err := ioutil.ReadFile("testdata/" + filename + ".xml")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func TestPostsAdd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/add?auth_token=user%3Atoken&description=Title&dt=2006-01-02+15%3A04%3A00+%2B0000+MST&extended=Description&replace=true&shared=true&tags=one&tags=two&tags=three&tags=four&toread=true&url=http%3A%2F%2Fexample.org",
		httpmock.NewStringResponder(200, readFixture("posts_ok")))

	tags := []string{"one", "two", "three", "four"}
	creation, _ := time.Parse(time.RFC822, time.RFC822)
	_, err := client.Posts.Add("http://example.org", "Title", "Description", tags, &creation, true, true, true)
	if err != nil {
		t.Error(err)
	}
}

func TestPostsDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/delete?auth_token=user%3Atoken&url=http%3A%2F%2Fexample.org",
		httpmock.NewStringResponder(200, readFixture("posts_ok")))

	_, err := client.Posts.Delete("http://example.org")
	if err != nil {
		t.Error(err)
	}
}

func TestPostsGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/get?auth_token=user%3Atoken&tags=webdev",
		httpmock.NewStringResponder(200, readFixture("posts_get")))

	tags := []string{"webdev"}
	posts, _, err := client.Posts.Get(tags, nil, "")
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 1 {
		t.Error("Retrieved wrong amount")
	}

	if strings.Compare(posts[0].URL, "http://www.howtocreate.co.uk/tutorials/texterise.php?dom=1") != 0 {
		t.Error("Retrieved wrong amount")
	}
}

func TestPostsGetTooManyTags(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tags := []string{"one", "two", "three", "four"}
	_, _, err := client.Posts.Get(tags, nil, "")
	if err == nil {
		t.Error(err)
	}
}

func TestPostsRecent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/recent?auth_token=user%3Atoken&count=5",
		httpmock.NewStringResponder(200, readFixture("posts_recent")))

	posts, _, err := client.Posts.Recent(nil, 5)
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 5 {
		t.Errorf("Retrieved wrong amount expected 5 got %d", len(posts))
	}
}

func TestPostsRecentHighCount(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	_, _, err := client.Posts.Recent(nil, 1000)
	if err.Error() != "count must be below 100" {
		t.Error(err)
	}

}
