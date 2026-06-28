package config_test

import (
	"testing"

	"github.com/Khaym03/REG/internal/config"
)

func TestBrowserConfig_Equal(t *testing.T) {

	tests := []struct {
		name string
		cfg  config.BrowserConfig
		want bool
	}{
		{
			name: "Equal fiels",
			cfg: config.BrowserConfig{
				Headless: false,
				Trace:    true,
			},
			want: true,
		},
		{
			name: "Diff fiels",
			cfg: config.BrowserConfig{
				Headless: true,
				Trace:    false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := config.BrowserConfig{
				Headless: false,
				Trace:    true,
			}

			got := bc.Equal(tt.cfg)

			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
