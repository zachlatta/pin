package pin

import (
	"net/http"
	"testing"

	"errors"
	"fmt"

	"github.com/jarcoal/httpmock"
)

var clientRequestCodes = []struct {
	in_code int
	out_err error
}{
	{http.StatusOK, nil},
	{http.StatusTooManyRequests, errors.New(http.StatusText(http.StatusTooManyRequests))},
	{http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized))},
}

func TestClientRateLimit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range clientRequestCodes {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "https://api.pinboard.in/v1/endpoint?auth_token=user%3Atoken",
			httpmock.NewStringResponder(tt.in_code, readFixture("posts_err")))

		req, err := client.NewRequest("endpoint", nil)
		if err != nil {
			t.Error(err)
		}
		_, err = client.Do(req, nil)

		//t.Error(err)

		if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.out_err) {
			t.Errorf("Missing HTTP handling for status %d got '%v' expected '%v'", tt.in_code, err, tt.out_err)
		}
	}
}
