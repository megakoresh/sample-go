package util

import (
	"log"
	"os"
)

var (
	Logger = log.New(os.Stderr, "SAMPLE-GO: ", log.LstdFlags|log.LUTC)
)

// GetString returns first argument or second if first is nil
func GetString(maybenil string, def string) string {
	if maybenil == "" {
		return def
	}
	return maybenil
}
