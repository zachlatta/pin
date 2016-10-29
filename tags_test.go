package pin

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestTagsGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/tags/get?auth_token=user%3Atoken",
		httpmock.NewStringResponder(200, readFixture("tags_get")))

	tags, _, err := client.Tags.Get()
	if err != nil {
		t.Error(err)
	}

	if len(tags) != 6 {
		t.Errorf("Wrong tags amount expected 6 got %d", len(tags))
	}

	if tags[2].Count != 3 {
		t.Errorf("Count value retrievel error expected 3 got %d", tags[2].Count)
	}
	if tags[2].Name != "radio" {
		t.Errorf("Name value retrievel error expected 'radio' got %s", tags[2].Name)
	}
}

func TestTagsRename(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/tags/rename?auth_token=user%3Atoken&new=new&old=old",
		httpmock.NewStringResponder(200, readFixture("ok")))

	_, err := client.Tags.Rename("new", "old")
	if err != nil {
		t.Error(err)
	}
}

func TestTagsDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/tags/delete?auth_token=user%3Atoken&tag=fooo",
		httpmock.NewStringResponder(200, readFixture("ok")))

	_, err := client.Tags.Delete("fooo")
	if err != nil {
		t.Error(err)
	}
}
