package mvd_parser

type Demo struct {
	Filename   string   `json:"filename"`
	Date       string   `json:"date"`
	Duration   string   `json:"duration"`
	ServerInfo string   `json:"serverinfo"`
	Players    []Player `json:"players"`
}

type Player struct {
	Name string `json:"name"`
	Team string `json:"team"`
}
