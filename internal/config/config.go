package config

import (
	"io"
)

type BrowserConfig struct {
	LoggerOut io.Writer
	Headless  bool `json:"headless,omitempty"`
	Trace     bool `json:"trace,omitempty"`
}

func (bc BrowserConfig) Equal(cfg BrowserConfig) bool {
	return bc.Headless == cfg.Headless &&
		bc.Trace == cfg.Trace
}
