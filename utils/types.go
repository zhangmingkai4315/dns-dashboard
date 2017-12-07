package utils

type Global struct {
	Type   string `ini:"type"`
	Server string `ini:"server"`
	Port   int    `ini:"port"`
	Source string `ini:"source"`
	Path   string `ini:"path"`
	Grok   string `ini:"grok"`
}

type DNS struct {
	Grok string `ini:"grok"`
}

type Config struct {
	Global
	DNS
}
