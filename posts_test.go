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
	time1  = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	time2  = time.Date(2009, time.December, 10, 23, 0, 0, 0, time.UTC)
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

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/add?auth_token=user%3Atoken&description=Title&dt=2009-11-10T23%3A00%3A00Z&extended=Description&replace=true&shared=true&tags=one&tags=two&tags=three&tags=four&toread=true&url=http%3A%2F%2Fexample.org",
		httpmock.NewStringResponder(200, readFixture("posts_ok")))

	tags := []string{"one", "two", "three", "four"}
	_, err := client.Posts.Add("http://example.org", "Title", "Description", tags, &time1, true, true, true)
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
		t.Errorf("Retrieved wrong amount - expected 1 got %s", len(posts))
	}

	if strings.Compare(posts[0].URL, "http://www.howtocreate.co.uk/tutorials/texterise.php?dom=1") != 0 {
		t.Error("Retrieved wrong results")
	}
}

var postsGetUrlTests = []struct {
	in_tags     []string
	in_creation *time.Time
	in_urlstr   string
	out_url     string
}{
	{[]string{}, nil, "", "https://api.pinboard.in/v1/posts/get?auth_token=user%3Atoken"},
	{[]string{"web", "dev"}, nil, "", "https://api.pinboard.in/v1/posts/get?auth_token=user%3Atoken&tags=web+dev"},
	{[]string{}, &time1, "", "https://api.pinboard.in/v1/posts/get?auth_token=user%3Atoken&dt=2009-11-10T23%3A00%3A00Z"},
	{[]string{}, nil, "http://example.org", "https://api.pinboard.in/v1/posts/get?auth_token=user%3Atoken&url=http%3A%2F%2Fexample.org"},
}

func TestPostsGetUrls(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range postsGetUrlTests {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", tt.out_url,
			httpmock.NewStringResponder(200, readFixture("posts_get")))
		_, _, err := client.Posts.Get(tt.in_tags, tt.in_creation, tt.in_urlstr)
		if err != nil {
			t.Error(err)
		}
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

func TestPostsAll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/all?auth_token=user%3Atoken&tags=webdev",
		httpmock.NewStringResponder(200, readFixture("posts_all")))

	tags := []string{"webdev"}
	posts, _, err := client.Posts.All(tags, 0, 0, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 2 {
		t.Errorf("Retrieved wrong amount - expected 2 got %s", len(posts))
	}

	if strings.Compare(posts[0].URL, "http://www.weather.com/") != 0 {
		t.Errorf("Retrieved wrong results (%s)", posts[0].URL)
	}
}

var postsAllUrlTests = []struct {
	in_tags    []string
	in_start   int
	in_results int
	in_fromdt  *time.Time
	in_todt    *time.Time
	out_url    string
}{
	{[]string{}, 0, 0, nil, nil, "https://api.pinboard.in/v1/posts/all?auth_token=user%3Atoken"},
	{[]string{"webdev"}, 0, 0, nil, nil, "https://api.pinboard.in/v1/posts/all?auth_token=user%3Atoken&tags=webdev"},
	{[]string{}, 10, 300, nil, nil, "https://api.pinboard.in/v1/posts/all?auth_token=user%3Atoken&results=300&start=10"},
	{[]string{}, 0, 0, &time1, &time2, "https://api.pinboard.in/v1/posts/all?auth_token=user%3Atoken&fromdt=2009-11-10T23%3A00%3A00Z&todt=2009-12-10T23%3A00%3A00Z"},
}

func TestPostsAllUrls(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range postsAllUrlTests {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", tt.out_url,
			httpmock.NewStringResponder(200, readFixture("posts_all")))
		_, _, err := client.Posts.All(tt.in_tags, tt.in_start, tt.in_results, tt.in_fromdt, tt.in_todt)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestPostsUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/update?auth_token=user%3Atoken",
		httpmock.NewStringResponder(200, readFixture("posts_update")))

	upd, _, err := client.Posts.LastTimeUpdated()
	if err != nil {
		t.Error(err)
	}

	if upd.Format(timeLayoutFull) != "2011-03-24T19:02:07Z" {
		t.Error("Wrong time recieved, expected 2011-03-24T19:02:07Z got %s", upd.Format(timeLayoutFull))
	}
}

func TestPostsDates(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/posts/dates?auth_token=user%3Atoken&tags=argentina",
		httpmock.NewStringResponder(200, readFixture("posts_dates")))

	tags := []string{"argentina"}
	dates, _, err := client.Posts.Dates(tags)
	if err != nil {
		t.Error(err)
	}

	if len(dates) != 8 {
		t.Errorf("Retrieved wrong amount - expected 8 got %d", len(dates))
	} else if dates[1].Count != 15 {
		t.Errorf("Retrieved wrong count - expected 15 got %d", dates[1].Count)
	}
}
