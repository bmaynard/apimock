package config

var Configuration = FileConfig{}

type FileConfig struct {
	Filesystem FilesystemConfig `mapstructure:"filesystem"`
}

type FilesystemConfig struct {
	Adapter string      `mapstructure:"adapter"`
	Meta    interface{} `mapstructure:"meta"`
}
