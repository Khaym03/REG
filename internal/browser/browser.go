package browser

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type BrowserConfig struct {
	LoggerOut io.Writer
}

func BuildBrowser(ctx context.Context, conf BrowserConfig) *rod.Browser {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	l := launcher.New().
		Headless(os.Getenv("REG_HEADLESS") == "1").
		Devtools(false).
		Leakless(false).
		UserDataDir(filepath.Join(rootDir, "rod_data"))

	browser := rod.New().
		Context(ctx).
		ControlURL(l.MustLaunch()).
		Trace(os.Getenv("REG_ROD_VERBOSE") == "1")

	if conf.LoggerOut != nil {
		l = l.Logger(conf.LoggerOut)
	}

	return browser.MustConnect()
}
