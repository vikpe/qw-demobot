package monitor

import (
	"fmt"
	"github.com/icza/backscanner"
	"os"
	"strings"
	"time"

	"github.com/vikpe/qw-demobot/internal/comms/topic"
)

type LogMonitor struct {
	isStopped        bool
	path             string
	onEvent          func(string, ...any)
	lastDemoFilename string
	logParser        *LogParser
}

func NewLogMonitor(path string, onEvent func(topic string, data ...any)) *LogMonitor {
	return &LogMonitor{
		isStopped:        false,
		path:             path,
		onEvent:          onEvent,
		lastDemoFilename: "",
		logParser:        NewLogParser(path),
	}
}

func (p *LogMonitor) Start(interval time.Duration) {
	p.isStopped = false
	ticker := time.NewTicker(interval)

	for ; true; <-ticker.C {
		if p.isStopped {
			return
		}

		p.CompareStates()
	}
}

func (p *LogMonitor) CompareStates() {
	currentDemoFilename := p.logParser.GetCurrentDemoFilename()
	p.logParser.Truncate()

	if currentDemoFilename != p.lastDemoFilename {
		p.onEvent(topic.DemoFilenameChanged, currentDemoFilename)
	}

	p.lastDemoFilename = currentDemoFilename
}

func (p *LogMonitor) Stop() {
	p.isStopped = true
}

// parse log file
const TOKEN_DEMO_START = "#bot#demo_start#"
const TOKEN_DEMO_END = "#bot#demo_stop#"

type LogParser struct {
	path string
}

func NewLogParser(path string) *LogParser {
	// create file if not exists
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("LogParser: %s", err))
	}
	defer file.Close()

	return &LogParser{
		path,
	}
}

func (l *LogParser) Truncate() {
	file, err := os.OpenFile(l.path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Sprintf("LogParser: %s", err))
	}
	defer file.Close()
}

func (l *LogParser) GetCurrentDemoFilename() string {
	file, err := os.Open(l.path)
	if err != nil {
		panic(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := backscanner.New(file, int(fileInfo.Size()))

	for {
		line, _, err := scanner.Line()

		if err != nil {
			return ""
		}

		if strings.HasPrefix(line, TOKEN_DEMO_END) {
			return ""
		} else if strings.HasPrefix(line, TOKEN_DEMO_START) {
			return strings.TrimPrefix(line, TOKEN_DEMO_START)
		}
	}

	return ""
}
