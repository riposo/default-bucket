package defaultbucket

import (
	"github.com/riposo/default-bucket/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
	"github.com/riposo/riposo/pkg/riposo"
)

func init() {
	plugin.Register(plugin.New(
		"default-bucket",
		map[string]interface{}{
			"description": "The default bucket is an alias for a personal bucket where collections are created implicitly.",
			"url":         "https://github.com/riposo/default-bucket",
		},
		func(rts *api.Routes) error {
			cfg := new(internal.Config)
			if err := riposo.ParseEnv(cfg); err != nil {
				return err
			}
			cfg.Mount(rts)
			return nil
		},
		nil,
	))
}
