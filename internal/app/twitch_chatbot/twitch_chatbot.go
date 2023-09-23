package twitch_chatbot

import (
	"fmt"
	"github.com/vikpe/qw-demobot/internal/comms/topic"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/vikpe/go-qwhub"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/qw-demobot/internal/app/twitch_chatbot/monitor"
	"github.com/vikpe/qw-demobot/internal/comms/commander"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq/message"
	hubTwitch "github.com/vikpe/qw-hub-api/pkg/twitch"
	chatbot "github.com/vikpe/twitch-chatbot"
)

func New(channelName, channelAccessToken, subscriberAddress, publisherAddress string) *chatbot.Chatbot {
	var pfmt = prettyfmt.New("chatbot", color.FgHiMagenta, "15:04:05", color.FgWhite)

	oauth_token := fmt.Sprintf("oauth:%s", channelAccessToken)
	bot := chatbot.NewChatbot(channelName, oauth_token, channelName, '!')

	// zmq messages
	subscriber := zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll)
	subscriber.OnMessage = func(message message.Message) {
		switch message.Topic {
		case topic.TwitchChatbotSay:
			bot.Say(message.Content.ToString())
		}
	}

	// announce when streamers go live
	streamsMonitor := monitor.NewStreamsMonitor(channelName, qwhub.NewClient().Streams, func(stream hubTwitch.Stream) {
		bot.Say(fmt.Sprintf("%s is now streaming @ %s - %s", stream.ClientName, stream.Url, stream.Title))
	})

	// bot events
	bot.OnConnected = func() {
		pfmt.Println("connected as", channelName)
	}

	bot.OnStarted = func() {
		pfmt.Println("started")
		go subscriber.Start()
		go streamsMonitor.Start(15 * time.Second)
	}

	bot.OnStopped = func(sig os.Signal) {
		subscriber.Stop()
		streamsMonitor.Stop()
		pfmt.Printfln("stopped (%s)", sig)
	}

	// channel commands
	cmder := commander.NewCommander(zeromq.NewPublisher(publisherAddress).SendMessage)

	bot.AddCommand("autotrack", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Autotrack()
	})

	bot.AddCommand("cfg_load", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.LoadConfig()
	})

	bot.AddCommand("cmd", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		if !chatbot.IsModerator(msg.User) {
			bot.Reply(msg, "cmd is a mod-only command.")
			return
		}

		cmder.Command(cmd.ArgsToString())
	})

	bot.AddCommand("commands", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		replyMessage := fmt.Sprintf(`available commands: %s`, bot.GetCommands(", "))
		bot.Reply(msg, replyMessage)
	})

	bot.AddCommand("help", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		replyMessage := fmt.Sprintf(`available commands: %s`, bot.GetCommands(", "))
		bot.Reply(msg, replyMessage)
	})

	bot.AddCommand("next", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		/*if !chatbot.IsModerator(msg.User) {
			bot.Reply(msg, "cmd is a mod-only command.")
			return
		}*/

		cmder.Command("demo_playlist_next")
	})

	bot.AddCommand("prev", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		/*if !chatbot.IsModerator(msg.User) {
			bot.Reply(msg, "cmd is a mod-only command.")
			return
		}*/

		cmder.Command("demo_playlist_prev")
	})

	bot.AddCommand("restart", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		if !chatbot.IsModerator(msg.User) {
			bot.Reply(msg, "restart is a mod-only command.")
			return
		}

		cmder.StopEzquake()
		time.AfterFunc(1250*time.Millisecond, func() {
			cmder.StopQuakeManager()
		})
	})

	bot.AddCommand("showscores", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Showscores()
	})

	bot.AddCommand("track", func(cmd chatbot.Command, msg twitch.PrivateMessage) {
		cmder.Track(cmd.ArgsToString())
	})

	return bot
}
