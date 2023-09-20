package monitor_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-demobot/internal/app/quake_manager/monitor"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq/mock"
	"os"
	"testing"
	"time"
)

// helpers
func appendToFile(path string, content string) {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if _, err := f.WriteString(content + "\n"); err != nil {
		panic(err)
	}
}

func TestLogMonitor(t *testing.T) {
	logPath := "test_logmonitor.log"
	eventCallback := mock.NewPublisherMock()
	logMonitor := monitor.NewLogMonitor(logPath, eventCallback.SendMessage)

	go logMonitor.Start(10 * time.Microsecond)
	appendToFile(logPath, "#bot#demo_start#duel_xantom_vs_bps.mvd")
	time.Sleep(time.Millisecond * 20)
	appendToFile(logPath, "#bot#demo_stop#")
	time.Sleep(time.Millisecond * 20)

	expectCalls := [][]any{
		{"demo.filename_changed", "duel_xantom_vs_bps.mvd"},
		{"demo.filename_changed", ""},
	}
	assert.Equal(t, expectCalls, eventCallback.SendMessageCalls)

	logMonitor.Stop()

	os.Remove(logPath) // cleanup
}

func TestLogParser(t *testing.T) {
	logPath := "test_logparser.log"
	logParser := monitor.NewLogParser(logPath)

	logParser.Truncate()
	assert.Equal(t, "", logParser.GetCurrentDemoFilename())

	appendToFile(logPath, "config loaded")
	assert.Equal(t, "", logParser.GetCurrentDemoFilename())

	appendToFile(logPath, "#bot#demo_start#duel_xantom_vs_bps.mvd")
	appendToFile(logPath, "match started")
	assert.Equal(t, "duel_xantom_vs_bps.mvd", logParser.GetCurrentDemoFilename())

	appendToFile(logPath, "1 minute left")
	assert.Equal(t, "duel_xantom_vs_bps.mvd", logParser.GetCurrentDemoFilename())

	appendToFile(logPath, "#bot#demo_stop#")
	assert.Equal(t, "", logParser.GetCurrentDemoFilename())

	appendToFile(logPath, "#bot#demo_start#duel_xantom_vs_bps.mvd")
	appendToFile(logPath, "#bot#demo_start#duel_xantom_vs_xterm.mvd")
	assert.Equal(t, "duel_xantom_vs_xterm.mvd", logParser.GetCurrentDemoFilename())

	os.Remove(logPath) // cleanup
}
