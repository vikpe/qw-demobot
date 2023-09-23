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
	isStopped bool
	path      string
	onEvent   func(string, ...any)
	lastDemo  string
	logParser *LogParser
}

func NewLogMonitor(path string, onEvent func(topic string, data ...any)) *LogMonitor {
	return &LogMonitor{
		isStopped: false,
		path:      path,
		onEvent:   onEvent,
		lastDemo:  "",
		logParser: NewLogParser(path),
	}
}

func (p *LogMonitor) Start(interval time.Duration) {
	p.isStopped = false
	ticker := time.NewTicker(interval)

	for ; true; <-ticker.C {
		if p.isStopped {
			return
		}

		p.compareStates()
	}
}

func (p *LogMonitor) compareStates() {
	currentDemo := p.logParser.GetDemo()
	hasStartedDemo := p.lastDemo == "" && currentDemo != ""
	hasStoppedDemo := p.lastDemo != "" && currentDemo == ""

	if hasStartedDemo {
		p.onEvent(topic.DemoStarted, currentDemo)
	} else if hasStoppedDemo {
		p.onEvent(topic.DemoStopped, p.lastDemo)
	}

	if currentDemo != p.lastDemo {
		p.onEvent(topic.DemoNameChanged, currentDemo)
	}

	p.lastDemo = currentDemo
}

func (p *LogMonitor) Stop() {
	p.isStopped = true
}

// parse log file
const TOKEN_DEMO_START = "#demo#start#"
const TOKEN_DEMO_END = "#demo#stop#"

type LogParser struct {
	path string
}

func NewLogParser(path string) *LogParser {
	return &LogParser{
		path,
	}
}

func (l *LogParser) touch() {
	file, err := os.OpenFile(l.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("LogParser: %s", err))
	}
	defer file.Close()
}

func (l *LogParser) GetDemo() string {
	l.touch()

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
