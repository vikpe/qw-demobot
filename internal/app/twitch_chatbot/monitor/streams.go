package monitor

import (
	"github.com/samber/lo"
	"github.com/vikpe/qw-demobot/internal/pkg/task"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

func NewStreamsMonitor(currentChannelName string, getStreams func() []twitch.Stream, onStreamStarted func(stream twitch.Stream)) *task.PeriodicalTask {
	var prevChannels []string

	onTick := func() {
		streams := getStreams()
		currentChannels := lo.Map(streams, func(stream twitch.Stream, _ int) string {
			return stream.Channel
		})

		if prevChannels == nil {
			prevChannels = currentChannels
			return
		}

		for _, stream := range streams {
			if stream.Channel == currentChannelName {
				continue
			}

			if !lo.Contains(prevChannels, stream.Channel) {
				onStreamStarted(stream)
			}
		}

		prevChannels = currentChannels
	}

	return task.NewPeriodicalTask(onTick)
}
