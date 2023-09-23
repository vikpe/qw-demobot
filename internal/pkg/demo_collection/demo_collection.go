package demo_collection

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/vikpe/qw-demobot/internal/pkg/futil"
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

func (d *DemoCollection) GetStage(demoName string) string {
	filePath, err := futil.FindFileAbsPath(d.Path, demoName)

	if err != nil || strings.Contains(filePath, "unsorted") {
		return ""
	}

	parts := strings.SplitN(demoName, "_", 3)
	return parts[1]
}

func (d *DemoCollection) GetMvdParserFilename(demoName string) string {
	return demoName + ".mvd.json"
}

func (d *DemoCollection) GetMvdParserInfo(demoName string) (mvd_parser.Demo, error) {
	infoFilename := d.GetMvdParserFilename(demoName)
	infoFilePath, err := futil.FindFileAbsPath(d.Path, infoFilename)

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

func (d *DemoCollection) GetFilename(demoName string) (string, error) {
	mvdFilename := demoName + ".mvd"
	if futil.DirHasFile(d.Path, mvdFilename) {
		return mvdFilename, nil
	}

	demFilename := demoName + ".dem"
	if futil.DirHasFile(d.Path, demFilename) {
		return demFilename, nil
	}

	return "", fmt.Errorf("Demo not found: %s", demoName)
}

func (d *DemoCollection) GetAbsPath(demoName string) (string, error) {
	mvdFilename := demoName + ".mvd"
	path, err := futil.FindFileAbsPath(d.Path, mvdFilename)
	if err == nil {
		return path, nil
	}

	demFilename := demoName + ".dem"
	path, err = futil.FindFileAbsPath(d.Path, demFilename)
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Demo not found: %s", demoName)
}

func (d *DemoCollection) GetEventInfo(demoName string) string {
	filename, err := d.GetFilename(demoName)

	if err != nil || strings.HasSuffix(filename, ".dem") {
		return ""
	}

	infoFilename := d.GetMvdParserFilename(demoName)
	infoFilePath, err := futil.FindFileAbsPath(d.Path, infoFilename)

	if err != nil || !strings.Contains(infoFilePath, "/tournaments") {
		return ""
	}

	relPath := strings.TrimPrefix(infoFilePath, d.Path+"/")
	dirs := strings.SplitN(relPath, "/", 3)
	eventName := strings.ReplaceAll(dirs[1], "_", " ")

	if strings.Count(demoName, "_") < 3 {
		return eventName
	}

	parts := strings.SplitN(demoName, "_", 3)
	stage := parts[1]
	return fmt.Sprintf("%s %s", eventName, stage)
}

func (d *DemoCollection) GetTitle(demoName string) string {
	info, err := d.GetMvdParserInfo(demoName)
	if err != nil {
		return demoName
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

func (d *DemoCollection) GetSha256(demoName string) (string, error) {
	fileName, err := d.GetFilename(demoName)

	if err != nil {
		return "", err
	}

	return futil.FindFileSha256(d.Path, fileName)
}
