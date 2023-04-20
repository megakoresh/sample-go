package send

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSend(t *testing.T) {
	cases := []struct {
		input, apikey string
		want          int
	}{
		{
			input:  `{"peruna": "potato", "porkkana": "carrot"}`,
			apikey: "correct",
			want:   0,
		},
		{
			input:  `{"peruna": "картофель", "porkkana": "морковь"}`,
			apikey: "correct",
			want:   0,
		},
		{
			input:  `{"peruna": "kartoffel", "porkkana": "karotte"}`,
			apikey: "wrong",
			want:   1,
		},
	}

	for i, c := range cases {
		t.Logf("Send test case %s", c.input)
		tmp, err := os.CreateTemp("", fmt.Sprintf("sample_test_send_%d.json", i))
		if err != nil {
			t.Fatalf("Could not create test file for case")
		}
		tmp.WriteString(c.input)
		tmp.Close()

		pastebinMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				t.Fatalf("Could not parse input of test case %s", err)
			}
			for k, v := range r.Form {
				if k == "api_dev_key" {
					if v[0] != c.apikey {
						t.Fatalf("Api dev key %s expected, was %s", c.apikey, v)
					}
					switch v[0] {
					case "correct":
						w.Write([]byte("OK!"))
					case "wrong":
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("BAD!"))
					default:
						t.Fatalf("Unrecognized api key %s", v[0])
					}
				}
			}
		}))

		res := Send([]string{fmt.Sprintf("-file=%s", tmp.Name()), fmt.Sprintf("-pbapikey=%s", c.apikey), fmt.Sprintf("-pburl=%s", pastebinMockServer.URL), "-destination=pastebin"})
		err = os.Remove(tmp.Name())
		if err != nil {
			t.Fatalf("Could not remove %s - %s", tmp.Name(), err)
		}
		pastebinMockServer.Close()
		if res != c.want {
			t.Fatalf("Test case failed %d != %d", res, c.want)
		}
	}
}
