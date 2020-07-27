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
		fa = NewLocalFileOptions()
		fa.SetOptionString("Path", m)
	case "s3":
		b, found := config.Configuration.Filesystem.Meta.(map[string]interface{})["bucket"].(string)

		if !found {
			l.Log.Fatal("\"bucket\" not found in configuration YAML.")
		}

		r, found := config.Configuration.Filesystem.Meta.(map[string]interface{})["region"].(string)

		if !found {
			l.Log.Fatal("\"region\" not found in configuration YAML.")
		}

		p, _ := config.Configuration.Filesystem.Meta.(map[string]interface{})["prefix"].(string)

		fa = NewS3FileOptions()
		fa.SetOptionString("Bucket", b)
		fa.SetOptionString("Region", r)
		fa.SetOptionString("Prefix", p)
	default:
		fa = NewLocalFileOptions()
		fa.SetOptionString("Path", "./mocks")
	}

	return fa
}
