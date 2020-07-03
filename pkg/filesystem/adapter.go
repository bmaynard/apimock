package filesystem

import (
	"github.com/bmaynard/apimock/pkg/config"
	l "github.com/bmaynard/apimock/pkg/utils/logger"
)

func GetAdapter() Filesystem {
	var fa Filesystem

	switch config.Configuration.Filesystem.Adapter {
	case "local":
		m, found := config.Configuration.Filesystem.Meta.(map[string]interface{})["mock_path"].(string)

		if !found {
			l.Log.Fatal("\"mock_path\" not found in configuration YAML.")
		}
		fa = NewLoclFileOptions()
		fa.SetOptionString("Path", m)
	default:
		fa = NewLoclFileOptions()
		fa.SetOptionString("Path", "./mocks")
	}

	return fa
}
