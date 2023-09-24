package mvdparser

import (
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"strings"
)

type Data struct {
	Filepath   string   `json:"filepath"`
	Date       string   `json:"date"`
	Duration   string   `json:"duration"`
	ServerInfo string   `json:"serverinfo"`
	Players    []Player `json:"players"`
}

type Player struct {
	Name        string `json:"name"`
	Team        string `json:"team"`
	TopColor    uint8  `json:"top_color"`
	BottomColor uint8  `json:"bottom_color"`
	Frags       int    `json:"frags"`
	Teamkills   uint   `json:"teamkills"`
	Deaths      uint   `json:"deaths"`
	Suicides    uint   `json:"suicides"`
	AvgPing     string `json:"avg_ping"`
}

func NewFromFile(filepath string) (*Data, error) {
	if !strings.HasSuffix(filepath, ".json") {
		return &Data{}, fmt.Errorf("File must be a json file: %s", filepath)
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &Data{}, err
	}

	var info Data
	err = json.Unmarshal(content, &info)

	return &info, err
}
