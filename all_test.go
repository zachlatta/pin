package pin

import "io/ioutil"

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
