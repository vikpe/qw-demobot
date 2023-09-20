package topic

// events
const (
	EzquakeStarted = "ezquake.started"
	EzquakeStopped = "ezquake.stopped"

	DemoFilenameChanged = "demo.filename_changed"
)

// commands
const (
	EzquakeCommand = "ezquake.command"
	EzquakeScript  = "ezquake.script"
	EzquakeStop    = "ezquake.stop"

	TwitchChatbotSay = "twitch_chatbot.say"

	QuakeManagerStop     = "quake_manager.stop"
	QuakeManagerEvaluate = "quake_manager.evaluate"
)
