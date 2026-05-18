package browser

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Khaym03/REG/internal/config"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func BuildBrowser(ctx context.Context, conf config.BrowserConfig) (*rod.Browser, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	l := launcher.New().
		Context(ctx).
		Headless(conf.Headless).
		Devtools(false).
		Leakless(false).
		UserDataDir(filepath.Join(rootDir, "rod_data"))

	controlURl, err := l.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().
		Context(ctx).
		ControlURL(controlURl).
		Trace(conf.Trace)

	if conf.LoggerOut != nil {
		l = l.Logger(conf.LoggerOut)
	}

	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	return browser, nil
}
