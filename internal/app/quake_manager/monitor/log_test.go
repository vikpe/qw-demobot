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
	os.OpenFile(logPath, os.O_RDWR|os.O_TRUNC, 0666)

	eventCallback := mock.NewPublisherMock()
	logMonitor := monitor.NewLogMonitor(logPath, eventCallback.SendMessage)

	go logMonitor.Start(20 * time.Microsecond)
	appendToFile(logPath, "#demo#start#duel_xantom_vs_bps")
	time.Sleep(time.Millisecond * 50)
	appendToFile(logPath, "#demo#stop#")
	time.Sleep(time.Millisecond * 50)

	expectCalls := [][]any{
		{"demo.started", "duel_xantom_vs_bps"},
		{"demo.name_changed", "duel_xantom_vs_bps"},
		{"demo.stopped", "duel_xantom_vs_bps"},
		{"demo.name_changed", ""},
	}
	assert.Equal(t, expectCalls, eventCallback.SendMessageCalls)

	logMonitor.Stop()

	os.Remove(logPath) // cleanup
}

func TestLogParser(t *testing.T) {
	logPath := "test_logparser.log"
	os.OpenFile(logPath, os.O_RDWR|os.O_TRUNC, 0666)

	logParser := monitor.NewLogParser(logPath)
	assert.Equal(t, "", logParser.GetDemo())

	appendToFile(logPath, "config loaded")
	assert.Equal(t, "", logParser.GetDemo())

	appendToFile(logPath, "#demo#start#duel_xantom_vs_bps")
	appendToFile(logPath, "match started")
	assert.Equal(t, "duel_xantom_vs_bps", logParser.GetDemo())

	appendToFile(logPath, "1 minute left")
	assert.Equal(t, "duel_xantom_vs_bps", logParser.GetDemo())

	appendToFile(logPath, "#demo#stop#")
	assert.Equal(t, "", logParser.GetDemo())

	appendToFile(logPath, "#demo#start#duel_xantom_vs_bps")
	appendToFile(logPath, "#demo#start#duel_xantom_vs_xterm")
	assert.Equal(t, "duel_xantom_vs_xterm", logParser.GetDemo())

	os.Remove(logPath) // cleanup
}
