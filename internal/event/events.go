package event

import (
	"github.com/mustafaturan/bus/v3"
	"github.com/mustafaturan/monoton/v2"
	"github.com/mustafaturan/monoton/v2/sequencer"
)

const (
	StatsTopic = "stats:result" // stats.Stats
)

const (
	WorkflowStartedTopic = "workflow.started"

	LogginTopic = "loggin"

	GuidesGatherTopic = "guides.gather"

	InventorySyncTopic = "inventory.sync"

	ReceptionTopic = "reception"

	LogoutTopic = "logout"

	WorkflowFinishedTopic = "workflow.finished"
)

type Topics struct {
	StatsResult           string `json:"stats_result"`
	WorkflowStartedTopic  string `json:"workflow_started_topic"`
	LogginTopic           string `json:"loggin_topic"`
	GuidesGatherTopic     string `json:"guides_gather_topic"`
	InventorySyncTopic    string `json:"inventory_sync_topic"`
	ReceptionTopic        string `json:"reception_topic"`
	LogoutTopic           string `json:"logout_topic"`
	WorkflowFinishedTopic string `json:"workflow_finished_topic"`
}

func StructTopics() Topics {
	return Topics{
		StatsResult:           StatsTopic,
		WorkflowStartedTopic:  WorkflowStartedTopic,
		LogginTopic:           LogginTopic,
		GuidesGatherTopic:     GuidesGatherTopic,
		InventorySyncTopic:    InventorySyncTopic,
		ReceptionTopic:        ReceptionTopic,
		LogoutTopic:           LogoutTopic,
		WorkflowFinishedTopic: WorkflowFinishedTopic,
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
	b.RegisterTopics(StatsTopic)

	return b
}
