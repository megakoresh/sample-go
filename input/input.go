package input

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/megakoresh/sample-go/util"
)

const (
	FmtJSON = "json"
	FmtXML  = "xml"

	FileStdin = "stdin"
)

type DataItem struct {
	Format string `json:"-" xml:"-"` // omit

	Peruna   string `json:"peruna" xml:"peruna"`
	Porkkana string `json:"porkkana" xml:"porkkana"`
}

func (di DataItem) String() string {
	return fmt.Sprintf("%s %s", di.Peruna, di.Porkkana)
}

func Parse(format string, file string) (*DataItem, error) {
	var (
		input *os.File
		di    DataItem
	)
	switch file {
	case FileStdin:
		util.Logger.Println("Reading from stdin")
		input = os.Stdin
	default:
		absPath, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		input, err = os.Open(absPath)
		defer input.Close()
		if err != nil {
			return nil, err
		}
	}

	switch format {
	case FmtJSON:
		di.Format = FmtJSON
		dec := json.NewDecoder(input)
		err := dec.Decode(&di)
		if err != nil && err != io.EOF {
			return nil, err
		}
		return &di, nil
	case FmtXML:
		di.Format = FmtXML
		dec := xml.NewDecoder(input)
		if err := dec.Decode(&di); err != nil {
			return nil, err
		}
		dec.Decode(&di)
	default:
		return nil, fmt.Errorf("unrecognized format: %s", format)
	}
	return &DataItem{}, nil
}
