package demo_collection

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/vikpe/qw-demobot/internal/pkg/ffind"
	"github.com/vikpe/qw-demobot/internal/pkg/mvd_parser"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qtext/qstring"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type DemoCollection struct {
	Path string
}

func New(path string) *DemoCollection {
	return &DemoCollection{
		path,
	}
}

func (d *DemoCollection) GetStage(filename string) string {
	filePath, err := ffind.FindFileAbsPath(d.Path, filename)

	if err != nil || strings.Contains(filePath, "unsorted") {
		return ""
	}

	parts := strings.SplitN(filename, "_", 3)
	return parts[1]
}

func (d *DemoCollection) GetInfoFilename(filename string) string {
	targetFilename := filename
	ext := filepath.Ext(filename)

	if ext != ".mvd" {
		targetFilename = strings.TrimSuffix(filename, ext) + ".mvd"
	}

	return targetFilename + ".json"
}

func (d *DemoCollection) GetMvdParserInfo(filename string) (mvd_parser.Demo, error) {
	infoFilename := d.GetInfoFilename(filename)
	infoFilePath, err := ffind.FindFileAbsPath(d.Path, infoFilename)

	if err != nil {
		return mvd_parser.Demo{}, err
	}

	content, err := ioutil.ReadFile(infoFilePath)

	if err != nil {
		return mvd_parser.Demo{}, err
	}

	var info mvd_parser.Demo
	json.Unmarshal(content, &info)

	return info, err
}

func (d *DemoCollection) GetEventInfo(filename string) string {
	infoFilename := filename + ".json"
	infoFilePath, err := ffind.FindFileAbsPath(d.Path, infoFilename)

	if err != nil || !strings.Contains(infoFilePath, "/tournaments") {
		return ""
	}

	relPath := strings.TrimPrefix(infoFilePath, d.Path+"/")
	dirs := strings.SplitN(relPath, "/", 3)
	eventName := strings.ReplaceAll(dirs[1], "_", " ")

	if strings.Count(filename, "_") < 3 {
		return eventName
	}

	parts := strings.SplitN(filename, "_", 3)
	stage := parts[1]
	return fmt.Sprintf("%s %s", eventName, stage)
}

func (d *DemoCollection) GetTitle(filename string) string {
	info, err := d.GetMvdParserInfo(filename)
	if err != nil {
		return filename
	}

	settings := qsettings.ParseString(info.ServerInfo)
	var clients []qclient.Client
	for _, player := range info.Players {
		clients = append(clients, qclient.Client{
			Name: qstring.QuakeString(player.Name),
			Team: qstring.QuakeString(player.Team),
		})
	}
	return qtitle.New(settings, clients)
}
