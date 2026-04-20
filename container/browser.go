package container

import (
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func BuildBrowser() *rod.Browser {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	l := launcher.New().
		Headless(os.Getenv("REG_HEADLESS") == "1").
		Devtools(false).
		Leakless(false).
		UserDataDir(filepath.Join(rootDir, "rod_data"))

	return rod.New().
		ControlURL(l.MustLaunch()).
		Trace(os.Getenv("REG_ROD_VERBOSE") == "1").
		MustConnect()
}
