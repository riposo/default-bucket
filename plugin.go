package main

import (
	"github.com/riposo/default-bucket/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
	"github.com/riposo/riposo/pkg/riposo"
)

var _ plugin.Factory = Plugin

// Plugin export definition.
func Plugin(rts *api.Routes) (plugin.Plugin, error) {
	cfg := new(internal.Config)
	if err := riposo.ParseEnv(cfg); err != nil {
		return nil, err
	}
	cfg.Mount(rts)

	return plugin.New(
		"default_bucket",
		map[string]interface{}{
			"description": "The default bucket is an alias for a personal bucket where collections are created implicitly.",
			"url":         "https://github.com/riposo/default-bucket",
		},
		nil,
	), nil
}

func main() {}
