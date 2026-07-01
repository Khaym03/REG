package event

import (
	"regexp"

	"github.com/mustafaturan/bus/v3"
	"github.com/mustafaturan/monoton/v2"
	"github.com/mustafaturan/monoton/v2/sequencer"
)

const ()

func Matcher(event string) string {
	return "^" + regexp.QuoteMeta(event) + "$"
}

type Topic string

const (
	WorkflowStarted Topic = "workflow.started"

	Login Topic = "loggin"

	Stats Topic = "stats:result" // stats.Stats

	GuidesGather Topic = "guides.gather"

	InventorySync Topic = "inventory.sync"

	Reception Topic = "reception"

	Logout Topic = "logout"

	WorkflowFinished Topic = "workflow.finished"

	BuildingBrowser  Topic = "bulding.browser"
	DestroyingBrowser Topic = "destroying.browser"
)


func NewBus() *bus.Bus {
	// configure id generator (it doesn't have to be monoton)
	node := uint64(1)
	initialTime := uint64(1577865600000) // set 2020-01-01 PST as initial time
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		panic(err)
	}

	// init an id generator
	var idGenerator bus.Next = m.Next

	// create a new bus instance
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		panic(err)
	}

	b.RegisterTopics(
		string(Stats),
		string(WorkflowStarted),
		string(BuildingBrowser),
		string(Login),
		string(GuidesGather),
		string(InventorySync),
		string(Reception),
		string(Logout),
		string(DestroyingBrowser),
		string(WorkflowFinished),
	)

	return b
}
