package mvd

type MvdParserResult struct {
	Demo     Demo     `json:"demo"`
	Hostname string   `json:"hostname"`
	Map      Map      `json:"map"`
	Settings Settings `json:"settings"`
	Players  []Player `json:"players"`
}

type Demo struct {
	Filename string `json:"filename"`
	Time     string `json:"time"`
}

type Map struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

type Player struct {
	NameSanatized string `json:"name_sanatized"`
	NameRaw       string `json:"name_raw"`
	TeamSanatized string `json:"team_sanatized"`
	TeamRaw       string `json:"team_raw"`
	TopColor      string `json:"top_color"`
	BottomColor   string `json:"bottom_color"`
	Frags         string `json:"frags"`
	Deaths        string `json:"deaths"`
	Kills         string `json:"kills"`
	Teamkills     string `json:"teamkills"`
	AvgPacketloss string `json:"avg_packetloss"`
	AvgPing       string `json:"avg_ping"`
}

type Settings struct {
	Gamedir       string `json:"gamedir"`
	Fraglimit     string `json:"fraglimit"`
	Timelimit     string `json:"timelimit"`
	Deathmatch    string `json:"deathmatch"`
	Maxfps        string `json:"maxfps"`
	Teamplay      string `json:"teamplay"`
	ZEXT          string `json:"z_ext"`
	Fpd           string `json:"fpd"`
	Maxclients    string `json:"maxclients"`
	Maxspectators string `json:"maxspectators"`
	Watervis      string `json:"watervis"`
	Version       string `json:"version"`
	Mod           string `json:"mod"`
}
