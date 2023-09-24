package quake_manager

import (
	"fmt"
	"github.com/vikpe/qw-demobot/internal/pkg/calc"
	"github.com/vikpe/qw-demobot/internal/pkg/demo/collection"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/vikpe/go-ezquake"
	"github.com/vikpe/prettyfmt"
	"github.com/vikpe/qw-demobot/internal/app/quake_manager/monitor"
	"github.com/vikpe/qw-demobot/internal/comms/commander"
	"github.com/vikpe/qw-demobot/internal/comms/topic"
	"github.com/vikpe/qw-demobot/internal/pkg/task"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq/message"
)

var pfmt = prettyfmt.New("quakemanager", color.FgHiCyan, "15:04:05", color.FgWhite)

type QuakeManager struct {
	controller     *ezquake.ClientController
	logMonitor     *monitor.LogMonitor
	processMonitor *monitor.ProcessMonitor
	evaluateTask   *task.PeriodicalTask
	subscriber     *zeromq.Subscriber
	publisher      *zeromq.Publisher
	commander      *commander.Commander
	demos          *collection.Collection
	stopChan       chan os.Signal
}

func New(
	ezquakeProcessUsername string,
	ezquakeBinPath string,
	ezquakeLogPath string,
	publisherAddress string,
	subscriberAddress string,
) *QuakeManager {
	controller := ezquake.NewClientController(ezquakeProcessUsername, ezquakeBinPath)
	publisher := zeromq.NewPublisher(publisherAddress)
	subscriber := zeromq.NewSubscriber(subscriberAddress, zeromq.TopicsAll)

	manager := QuakeManager{
		controller:     controller,
		processMonitor: monitor.NewProcessMonitor(controller.Process.IsStarted, publisher.SendMessage),
		logMonitor:     monitor.NewLogMonitor(ezquakeLogPath, publisher.SendMessage),
		evaluateTask:   task.NewPeriodicalTask(func() { publisher.SendMessage(topic.QuakeManagerEvaluate) }),
		publisher:      publisher,
		subscriber:     subscriber,
		commander:      commander.New(publisher.SendMessage),
		demos:          collection.New("/home/vikpe/games/demoquake/qw/demos"),
	}
	subscriber.OnMessage = manager.OnMessage

	return &manager
}

func (m *QuakeManager) Start() {
	m.stopChan = make(chan os.Signal, 1)
	signal.Notify(m.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		// event listeners
		go m.subscriber.Start()
		zeromq.WaitForConnection()

		// event dispatchers
		go m.processMonitor.Start(3 * time.Second)
		go m.logMonitor.Start(2 * time.Second)

		if m.controller.Process.IsStarted() {
			go m.evaluateTask.Start(10 * time.Second)
		}
	}()
	<-m.stopChan
}

func (m *QuakeManager) Stop() {
	if m.stopChan == nil {
		return
	}
	m.stopChan <- syscall.SIGINT
	time.Sleep(30 * time.Millisecond)
}

func (m *QuakeManager) OnMessage(msg message.Message) {
	handlers := map[string]message.Handler{
		// commands
		topic.EzquakeCommand:   m.OnEzquakeCommand,
		topic.EzquakeScript:    m.OnEzquakeScript,
		topic.EzquakeStop:      m.OnStopEzquake,
		topic.QuakeManagerStop: m.OnStopQuakeManager,

		// ezquake events
		topic.EzquakeStarted: m.OnEzquakeStarted,
		topic.EzquakeStopped: m.OnEzquakeStopped,

		// demo events
		topic.DemoNameChanged: m.OnDemoNameChanged,
	}

	if handler, ok := handlers[msg.Topic]; ok {
		handler(msg)
	}
}

func (m *QuakeManager) OnStopQuakeManager(msg message.Message) {
	m.Stop()
}

func (m *QuakeManager) OnEzquakeCommand(msg message.Message) {
	m.controller.Command(msg.Content.ToString())
}

func (m *QuakeManager) OnEzquakeScript(msg message.Message) {
	script := msg.Content.ToString()

	switch script {
	case "load_config":
		m.commander.Command("cfg_load")
	case "showscores":
		m.commander.Command("+showscores")
		time.AfterFunc(8*time.Second, func() { m.commander.Command("-showscores") })
	}
}

func (m *QuakeManager) OnEzquakeStarted(msg message.Message) {
	pfmt.Println("OnEzquakeStarted")
	go m.evaluateTask.Start(10 * time.Second)
	time.AfterFunc(6*time.Second, func() { m.commander.Command("toggleconsole") })
}

func (m *QuakeManager) OnStopEzquake(msg message.Message) {
	pfmt.Println("OnStopEzquake")
	m.controller.Process.Stop(syscall.SIGTERM)

	time.AfterFunc(1*time.Second, func() {
		if m.controller.Process.IsStarted() {
			m.controller.Process.Stop(syscall.SIGKILL)
		}
	})
}

func (m *QuakeManager) OnEzquakeStopped(msg message.Message) {
	pfmt.Println("OnEzquakeStopped")
	m.evaluateTask.Stop()
}

func (m *QuakeManager) OnDemoNameChanged(msg message.Message) {
	demoName := msg.Content.ToString()
	pfmt.Println("OnDemoNameChanged", demoName)

	inGameTitle := demoName
	twitchTitle := "24/7 QuakeWorld demos"

	if len(demoName) > 0 {
		eventInfo := m.demos.GetEventInfo(demoName)
		demoTitle := m.demos.GetTitle(demoName)
		twitchTitle = demoTitle

		if len(eventInfo) > 0 {
			inGameTitle = eventInfo
			twitchTitle = fmt.Sprintf("%s / %s", eventInfo, demoTitle)
		}
	}

	if len(inGameTitle) > 0 {
		m.commander.Commandf("hud_static_text_scale %f", calc.StaticTextScale(inGameTitle))
	}

	fmt.Println("\ninGameTitle", inGameTitle)
	fmt.Println("twitchTitle", twitchTitle)
	sha256, err := m.demos.GetSha256(demoName)
	fmt.Println("sha256", sha256, err)

	m.commander.Commandf("bot_set_statictext %s", inGameTitle)
	m.publisher.SendMessage(topic.TwitchChannelSetTitle, twitchTitle)
}
