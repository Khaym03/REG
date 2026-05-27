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

const (
	WorkflowStarted = "workflow.started"

	Login = "loggin"

	Stats = "stats:result" // stats.Stats

	GuidesGather = "guides.gather"

	InventorySync = "inventory.sync"

	Reception = "reception"

	Logout = "logout"

	WorkflowFinished = "workflow.finished"

	BuildingBrowser   = "bulding.browser"
	DestroyingBrowser = "destroying.browser"
)

type Topics struct {
	Stats             string `json:"stats_result"`
	WorkflowStarted   string `json:"workflow_started"`
	Login             string `json:"login"`
	GuidesGather      string `json:"guides_gather"`
	InventorySync     string `json:"inventory_sync"`
	Reception         string `json:"reception"`
	Logout            string `json:"logout"`
	WorkflowFinished  string `json:"workflow_finished"`
	BuildingBrowser   string `json:"building_browser"`
	DestroyingBrowser string `json:"destroying_browser"`
}

func All() Topics {
	return Topics{
		Stats:             Stats,
		WorkflowStarted:   WorkflowStarted,
		BuildingBrowser:   BuildingBrowser,
		Login:             Login,
		GuidesGather:      GuidesGather,
		InventorySync:     InventorySync,
		Reception:         Reception,
		Logout:            Logout,
		DestroyingBrowser: DestroyingBrowser,
		WorkflowFinished:  WorkflowFinished,
	}
}

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

	// maybe register topics in here
	b.RegisterTopics(
		Stats,
		WorkflowStarted,
		BuildingBrowser,
		Login,
		GuidesGather,
		InventorySync,
		Reception,
		Logout,
		DestroyingBrowser,
		WorkflowFinished,
	)

	return b
}
