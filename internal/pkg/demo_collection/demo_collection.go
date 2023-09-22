package demo_collection

import (
	"github.com/goccy/go-json"
	"github.com/vikpe/qw-demobot/internal/pkg/ffind"
	"github.com/vikpe/qw-demobot/internal/pkg/mvd_parser"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qtext/qstring"
	"io/ioutil"
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
	parts := strings.SplitN(filename, "_", 3)
	return parts[1]
}

func (d *DemoCollection) GetInfo(filename string) (mvd_parser.Demo, error) {
	infoFilename := filename + ".json"
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

func (d *DemoCollection) GetEvent(filename string) string {
	infoFilename := filename + ".json"
	infoFilePath, err := ffind.FindFileAbsPath(d.Path, infoFilename)

	if err != nil {
		return "unknown"
	}

	relPath := strings.TrimPrefix(infoFilePath, d.Path+"/")
	dirs := strings.SplitN(relPath, "/", 3)
	eventName := strings.ReplaceAll(dirs[1], "_", " ")
	return eventName
}

func (d *DemoCollection) GetTitle(filename string) string {
	info, err := d.GetInfo(filename)
	if err != nil {
		return "unknown"
	}

	settings := qsettings.ParseString(info.Settings.ServerInfo)
	var clients []qclient.Client
	for _, player := range info.Players {
		clients = append(clients, qclient.Client{
			Name: qstring.QuakeString(player.Name),
			Team: qstring.QuakeString(player.Team),
		})
	}
	return qtitle.New(settings, clients)
}
