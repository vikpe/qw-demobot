package export

import (
	"fmt"
	"github.com/vikpe/qw-demobot/internal/pkg/demo/info"
	"github.com/vikpe/qw-demobot/internal/pkg/futil"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"path/filepath"
	"strings"
)

type Export struct {
	Sha256    string `json:"sha256"`
	Filepath  string `json:"filepath"`
	Filename  string `json:"filename"`
	Timestamp string `json:"timestamp"`
	Mode      string `json:"mode"`
	Map       string `json:"map"`
	Title     string `json:"title"`
	Event     string `json:"event"`
	Round     string `json:"round"`
	BestOf    int    `json:"best_of"`
	MapNumber int    `json:"map_number"`
}

func NewFromInfoPath(infoPath string) (Export, error) {
	info, err := info.NewFromMvdparserDataFile(infoPath)
	if err != nil {
		return Export{}, err
	}

	demoPath := strings.ReplaceAll(infoPath, ".mvd.json", ".mvd")

	sha256, err := futil.FileSha256(demoPath)
	if err != nil {
		fmt.Println("unable to get sha256 for", demoPath)
	}

	mode, _ := qmode.Parse(info.Settings)

	return Export{
		Sha256:    sha256,
		Filepath:  info.Filepath,
		Filename:  filepath.Base(info.Filepath),
		Timestamp: info.Timestamp.Format("2006-01-02 15:04"),
		Mode:      string(mode),
		Map:       info.Settings.Get("map", ""),
		Title:     info.GetTitle(),
		Event:     "todo",
		Round:     "todo",
		BestOf:    1,
		MapNumber: 1,
	}, nil
}
