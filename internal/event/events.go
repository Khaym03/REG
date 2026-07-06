package event

import (
	"regexp"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const ()

func Matcher(event string) string {
	return "^" + regexp.QuoteMeta(event) + "$"
}

type Empty struct{}

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

	BuildingBrowser   Topic = "bulding.browser"
	DestroyingBrowser Topic = "destroying.browser"
)

type Bus interface {
	Emit(name string, data any)
}

var AviableTopis = [...]string{
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
}

var _ Bus = (*bus)(nil)

type bus struct {
	event *application.EventManager
}

// Emit implements [Bus].
func (b *bus) Emit(name string, data any) {
	b.event.Emit(name, data)
}

func NewBus(em *application.EventManager) Bus {
	return &bus{
		event: em,
	}
}

var _ Bus = (*fakeBus)(nil)

type fakeBus struct{}

func NewFakeBus() *fakeBus {
	return new(fakeBus)
}

// Emit implements [Bus].
func (f *fakeBus) Emit(_ string, _ any) {}
