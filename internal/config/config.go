package config

import (
	"io"
	"os"
)

func IsDev() bool {
	return os.Getenv("APP_ENV") == "dev"
}

func IsProd() bool {
	return os.Getenv("APP_ENV") == "prod"
}

type BrowserConfig struct {
	LoggerOut io.Writer
	Headless  bool `json:"headless,omitempty"`
	Trace     bool `json:"trace,omitempty"`
}

func (bc BrowserConfig) Equal(cfg BrowserConfig) bool {
	return bc.Headless == cfg.Headless &&
		bc.Trace == cfg.Trace
}

func BrowserConfFromENV() BrowserConfig {
	return BrowserConfig{
		Headless: os.Getenv("REG_HEADLESS") == "1",
		Trace:    os.Getenv("REG_ROD_VERBOSE") == "1",
	}
}
