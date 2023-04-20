package send

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/megakoresh/sample-go/input"
	"github.com/megakoresh/sample-go/util"
)

const (
	destPastebin = "pastebin"
	pastebinUrl  = "https://pastebin.com/api/api_post.php"
)

var (
	fs *flag.FlagSet

	destination  string
	pbApiKey     string
	pbPasteTitle string
	help         bool
	format       string
	file         string
	pbUrl        string
)

func init() {
	fs = flag.NewFlagSet("send", flag.CommandLine.ErrorHandling())
	fs.StringVar(&file, "file", input.FileStdin, "File for processing")
	fs.StringVar(&format, "format", input.FmtJSON, "Input file format")
	fs.StringVar(&destination, "destination", destPastebin, "Where to send the data")
	fs.StringVar(&pbUrl, "pburl", pastebinUrl, "Pastebin url")
	fs.StringVar(&pbApiKey, "pbapikey", util.GetString(os.Getenv("PB_API_KEY"), ""), "If destination is pastebin, then api key for it (required)")
	fs.StringVar(&pbPasteTitle, "pbpastetitle", "swiggity.json", "Title of your pastebin paste if destination is pastebin")
	fs.BoolVar(&help, "help", false, "Print this help")
}

// it's usually best to return to main goroutine for any kind of user output. Don't print output from spawned goroutines
func doSend(di *input.DataItem, rc chan<- *http.Response, errChan chan<- error) {
	switch destination {
	case destPastebin:
		if pbApiKey == "" {
			errChan <- fmt.Errorf("no pastebin api key supplied, can't send (╯°□°）╯︵ ┻━┻")
			return
		}
		if pbUrl == "" {
			errChan <- fmt.Errorf("no pastebin url supplied, can't send (╯°□°）╯︵ ┻━┻")
			return
		}
		indentedJson, err := json.MarshalIndent(di, "", "  ")
		if err != nil {
			errChan <- err
			return
		}
		form := url.Values{}
		form.Add("api_dev_key", pbApiKey)
		form.Add("api_paste_code", string(indentedJson))
		form.Add("api_paste_name", pbPasteTitle)
		form.Add("api_paste_format", di.Format)
		form.Add("api_paste_expire", "1W")
		form.Add("api_option", "paste")
		res, err := http.PostForm(pbUrl, form)
		if err != nil {
			errChan <- err
			return
		}
		rc <- res
		return
	default:
		errChan <- fmt.Errorf("unsupported destination: %s", destination)
		return
	}
}

func Send(args []string) int {
	util.Logger.Printf("Sending input to %s API", destination)
	fs.Parse(args)
	if help {
		fs.Usage()
		os.Exit(0)
	}
	if destination == "" {
		util.Logger.Printf("no url specified, cannot send")
		return 1
	}

	di, err := input.Parse(format, file)
	if err != nil {
		util.Logger.Printf("Error while parsing input: %v", err)
		return 1
	}

	var (
		resChan = make(chan *http.Response)
		errChan = make(chan error)
	)

	go doSend(di, resChan, errChan)

	for {
		select {
		case r := <-resChan:
			defer r.Body.Close()
			util.Logger.Println("Received response from api server")
			_, err := io.Copy(os.Stdout, r.Body) // this is used if you have potentially big streaming response that you want to pipe directly to
			os.Stdout.WriteString("\n")          // terminator to avoid shell pollution
			if err != nil {
				util.Logger.Printf("Error while piping api response to stdout: %v", err)
				return 1
			}
			if r.StatusCode >= 400 {
				util.Logger.Printf("API server returned error code: %d", r.StatusCode)
				return 1
			}
			return 0
		case e := <-errChan:
			util.Logger.Println("Error while sending data to api server")
			util.Logger.Println(e)
			return 1
		}
	}
}
