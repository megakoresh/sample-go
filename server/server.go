package server

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/megakoresh/sample-go/util"
)

const (
	moodHappy = "happy"
	moodSad   = "sad"
)

var (
	fs *flag.FlagSet
)

func init() {
	fs = flag.NewFlagSet("server", flag.CommandLine.ErrorHandling())
}

type apiHandler struct {
	mood string
}

func (api apiHandler) MoodHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Plz use GET")
		return
	}
	switch api.mood {
	case moodHappy:
		fmt.Fprintf(w, "I am happy poikka!")
	case moodSad:
		fmt.Fprintf(w, "I am very surullinen boi :(")
	}
}

func (api apiHandler) Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/mood", api.MoodHandler)
	return router
}

func createServer(mood string) *http.ServeMux {
	api := apiHandler{mood: mood}
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.Router()))
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	return mux
}

func Serve(args []string) *http.ServeMux {
	// optionally use environment var
	mood := fs.String("mood", util.GetString(os.Getenv("SAMPLE_MOOD"), moodHappy), "The mood server is in. Controls what the /api/mood will return")
	port := fs.Int("port", 8090, "Port to listen on")
	fs.Parse(args)

	mux := createServer(*mood)
	util.Logger.Printf("Listening on :%d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)

	return mux
}
