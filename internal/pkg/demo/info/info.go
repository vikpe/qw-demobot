package info

import (
	"github.com/vikpe/qw-demobot/internal/pkg/demo/mvdparser"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qtitle"
	"github.com/vikpe/serverstat/qtext/qstring"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	Filepath  string             `json:"filepath"`
	Timestamp time.Time          `json:"timestamp"`
	Settings  qsettings.Settings `json:"settings"`
	Players   []qclient.Client   `json:"players"`
}

func NewFromMvdparserDataFile(filepath string) (Info, error) {
	data, err := mvdparser.NewFromFile(filepath)
	if err != nil {
		return Info{}, err
	}

	return NewFromMvdparserData(*data), nil
}

func NewFromMvdparserData(data mvdparser.Data) Info {
	// timestamp
	dateParts := strings.SplitN(data.Date, " ", 2)
	timestamp, err := time.Parse("2006-1-2", dateParts[0])
	if err != nil {
		timestamp = time.Time{}
	}

	// players
	var players []qclient.Client

	for _, player := range data.Players {
		if isVoidPlayer(player) {
			continue
		}

		players = append(players, qclient.Client{
			Name: qstring.QuakeString(strings.TrimSpace(player.Name)),
			Team: qstring.QuakeString(strings.TrimSpace(player.Team)),
		})
	}

	// settings
	settings := qsettings.ParseString(data.ServerInfo)
	settings["maxclients"] = strconv.Itoa(calcMaxClients(settings, players))

	return Info{
		Filepath:  data.Filepath,
		Timestamp: timestamp,
		Settings:  settings,
		Players:   players,
	}
}

func (i *Info) GetTitle() string {
	settings := i.Settings
	settings["matchtag"] = ""
	return qtitle.New(settings, i.Players)
}

func (i *Info) GetMode() qmode.Mode {
	mode, _ := qmode.Parse(i.Settings)
	return mode
}

func calcMaxClients(settings qsettings.Settings, players []qclient.Client) int {
	playerCount := len(players)

	isTeamplay := settings.GetInt("teamplay", 0) > 0
	if isTeamplay {
		return playerCount
	}

	mode, _ := qmode.Parse(settings)
	if mode.IsFfa() && 2 == playerCount {
		return 2
	}

	return settings.GetInt("maxclients", playerCount)
}

func isVoidPlayer(player mvdparser.Player) bool {
	return player.Frags == 0 && player.Teamkills == 0 && player.Deaths == 0 && player.Suicides == 0
}
