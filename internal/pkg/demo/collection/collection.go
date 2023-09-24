package collection

import (
	"fmt"
	"github.com/vikpe/qw-demobot/internal/pkg/demo/export"
	"github.com/vikpe/qw-demobot/internal/pkg/demo/info"
	"github.com/vikpe/qw-demobot/internal/pkg/futil"
	"sort"
	"strings"
)

type Collection struct {
	Path string
}

func New(path string) *Collection {
	return &Collection{
		path,
	}
}

func (d *Collection) GetStage(demoName string) string {
	filePath, err := futil.FindFilepath(d.Path, demoName)

	if err != nil || strings.Contains(filePath, "unsorted") {
		return ""
	}

	parts := strings.SplitN(demoName, "_", 3)
	return parts[1]
}

func (d *Collection) GetMvdParserFilename(demoName string) string {
	return demoName + ".mvd.json"
}

func (d *Collection) GetInfo(demoName string) (info.Info, error) {
	infoFilename := d.GetMvdParserFilename(demoName)
	infoFilePath, err := futil.FindFilepath(d.Path, infoFilename)

	if err != nil {
		return info.Info{}, err
	}

	return info.NewFromMvdparserDataFile(infoFilePath)
}

func (d *Collection) GetFilename(demoName string) (string, error) {
	mvdFilename := demoName + ".mvd"
	if futil.DirHasFile(d.Path, mvdFilename) {
		return mvdFilename, nil
	}

	demFilename := demoName + ".dem"
	if futil.DirHasFile(d.Path, demFilename) {
		return demFilename, nil
	}

	return "", fmt.Errorf("Data not found: %s", demoName)
}

func (d *Collection) GetAbsPath(demoName string) (string, error) {
	mvdFilename := demoName + ".mvd"
	path, err := futil.FindFilepath(d.Path, mvdFilename)
	if err == nil {
		return path, nil
	}

	demFilename := demoName + ".dem"
	path, err = futil.FindFilepath(d.Path, demFilename)
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Data not found: %s", demoName)
}

func (d *Collection) GetEventInfo(demoName string) string {
	filename, err := d.GetFilename(demoName)

	if err != nil || strings.HasSuffix(filename, ".dem") {
		return ""
	}

	infoFilename := d.GetMvdParserFilename(demoName)
	infoFilePath, err := futil.FindFilepath(d.Path, infoFilename)

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

func (d *Collection) GetTitle(demoName string) string {
	info, err := d.GetInfo(demoName)
	if err != nil {
		return demoName
	}

	return info.GetTitle()
}

func (d *Collection) GetSha256(demoName string) (string, error) {
	fileName, err := d.GetFilename(demoName)

	if err != nil {
		return "", err
	}

	return futil.FindFileSha256(d.Path, fileName)
}

func (d *Collection) Export(itemCallback func(string, export.Export) export.Export) []export.Export {
	infoPaths := futil.FindFilepathsByExtension(d.Path, ".mvd.json")
	sort.Strings(infoPaths)

	result := make([]export.Export, 0)

	// collect demos
	for _, infoPath := range infoPaths {
		demoExport, err := export.NewFromInfoPath(infoPath)

		if err != nil {
			fmt.Println("unable to export", infoPath)
			continue
		}

		demoPath := strings.ReplaceAll(infoPath, ".mvd.json", ".json")
		result = append(result, itemCallback(demoPath, demoExport))
	}

	return result
}
