package event_test

import (
	"context"
	"testing"

	"github.com/Khaym03/REG/internal/event"
	"github.com/mustafaturan/bus/v3"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler(t *testing.T) {
	buss := event.NewBus()

	expect := "hello"
	actual := ""
	statsHandler := bus.Handler{
		Handle: func(_ context.Context, e bus.Event) {
			actual = e.Data.(string)
		},
	}

	buss.RegisterHandler(event.StatsTopic, statsHandler)

	buss.Emit(t.Context(), event.StatsTopic, "hello")

	assert.Equal(t, expect, actual)
}
