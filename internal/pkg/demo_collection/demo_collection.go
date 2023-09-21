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

func (d *DemoCollection) GetEvent(filename string) (string, error) {
	infoFilename := filename + ".json"
	infoFilePath, err := ffind.FindFileAbsPath(d.Path, infoFilename)

	if err != nil {
		return "", err
	}

	relPath := strings.TrimPrefix(infoFilePath, d.Path+"/")
	dirs := strings.SplitN(relPath, "/", 3)
	eventName := strings.ReplaceAll(dirs[1], "_", " ")
	return eventName, nil
}

func (d *DemoCollection) GetTitle(filename string) (string, error) {
	info, err := d.GetInfo(filename)
	if err != nil {
		return "", err
	}

	settings := qsettings.ParseString(info.Settings.ServerInfo)
	var clients []qclient.Client
	for _, player := range info.Players {
		clients = append(clients, qclient.Client{
			Name: qstring.QuakeString(player.Name),
			Team: qstring.QuakeString(player.Team),
		})
	}
	title := qtitle.New(settings, clients)

	return title, nil
}
