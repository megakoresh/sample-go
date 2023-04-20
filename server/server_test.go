package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateServer(t *testing.T) {
	cases := []struct{ Name, m, want string }{
		{
			Name: "Sad case",
			m:    "sad",
			want: "I am very surullinen boi :(",
		}, {
			Name: "Happy case",
			m:    "happy",
			want: "I am happy poikka!",
		},
	}

	for _, c := range cases {
		t.Logf("Running %s", c.Name)
		s := createServer(c.m)
		testS := httptest.NewServer(s)
		defer testS.Close()
		res, err := http.Get(fmt.Sprintf("%v/api/mood", testS.URL))
		if err != nil {
			t.Fatalf("Test case failed with %s", err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("Test case failed reading response body with %s", err)
		}
		if c.want != string(body) {
			t.Fatalf("Test case failed %s != %s", c.want, body)
		}
	}
}
