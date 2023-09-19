package commander_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-demobot/internal/comms/commander"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq/mock"
)

func TestCommander_Autotrack(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Autotrack()

	expectedCalls := [][]any{{"ezquake.command", "bot_track"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Command(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Command("console")

	expectedCalls := [][]any{{"ezquake.command", "console"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Commandf(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Commandf("say %s", "foo")

	expectedCalls := [][]any{{"ezquake.command", "say foo"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Evaluate(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Evaluate()

	expectedCalls := [][]any{{"quake_manager.evaluate"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_LoadConfig(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.LoadConfig()

	expectedCalls := [][]any{{"ezquake.script", "load_config"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Showscores(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Showscores()

	expectedCalls := [][]any{{"ezquake.script", "showscores"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_StopEzquake(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.StopEzquake()

	expectedCalls := [][]any{{"ezquake.stop"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}

func TestCommander_Track(t *testing.T) {
	publisher := mock.NewPublisherMock()
	cmder := commander.NewCommander(publisher.SendMessage)
	cmder.Track("xantom")

	expectedCalls := [][]any{{"ezquake.command", "bot_track xantom"}}
	assert.Equal(t, expectedCalls, publisher.SendMessageCalls)
}
