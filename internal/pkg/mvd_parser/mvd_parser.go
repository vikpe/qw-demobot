package mvd_parser

type Demo struct {
	Filename string   `json:"filename"`
	Duration string   `json:"duration"`
	Settings Settings `json:"settings"`
	Players  []Player `json:"players"`
}

type Player struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

type Settings struct {
	ServerInfo string `json:"serverinfo"`
	Gamedir    string `json:"gamedir"`
	Fraglimit  string `json:"fraglimit"`
	Timelimit  string `json:"timelimit"`
	Deathmatch string `json:"deathmatch"`
	Teamplay   string `json:"teamplay"`
	Maxclients string `json:"maxclients"`
}
