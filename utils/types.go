package utils

type Global struct {
	Type   string `ini:"type"`
	Server string `ini:"server"`
	Port   int    `ini:"port"`
	Source string `ini:"source"`
	Path   string `ini:"path"`
}

type Config struct {
	Global
}
