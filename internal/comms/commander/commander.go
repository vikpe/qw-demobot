// Package commander sends commands
package commander

import (
	"fmt"

	"github.com/vikpe/qw-demobot/internal/comms/topic"
)

type Commander struct {
	sendMessage func(string, ...any)
}

func NewCommander(sendMessage func(topic string, data ...any)) *Commander {
	return &Commander{
		sendMessage: sendMessage,
	}
}

func (c Commander) Autotrack() {
	c.Commandf("bot_track")
}

func (c Commander) Command(cmd string) {
	c.sendMessage(topic.EzquakeCommand, cmd)
}

func (c Commander) Commandf(format string, args ...any) {
	c.Command(fmt.Sprintf(format, args...))
}

func (c Commander) Evaluate() {
	c.sendMessage(topic.QuakeManagerEvaluate)
}

func (c Commander) LoadConfig() {
	c.sendMessage(topic.EzquakeScript, "load_config")
}

func (c Commander) StopEzquake() {
	c.sendMessage(topic.EzquakeStop)
}

func (c Commander) StopQuakeManager() {
	c.sendMessage(topic.QuakeManagerStop)
}

func (c Commander) Track(target string) {
	c.sendMessage(topic.EzquakeCommand, fmt.Sprintf("bot_track %s", target))
}
