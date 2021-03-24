package defaultbucket

import (
	"github.com/riposo/default-bucket/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
	"github.com/riposo/riposo/pkg/riposo"
)

func init() {
	plugin.Register("default_bucket", func(rts *api.Routes) (plugin.Plugin, error) {
		cfg := new(internal.Config)
		if err := riposo.ParseEnv(cfg); err != nil {
			return nil, err
		}
		cfg.Mount(rts)

		return pin{
			"description": "The default bucket is an alias for a personal bucket where collections are created implicitly.",
			"url":         "https://github.com/riposo/default-bucket",
		}, nil
	})
}

type pin map[string]interface{}

func (p pin) Meta() map[string]interface{} { return map[string]interface{}(p) }
func (pin) Close() error                   { return nil }
