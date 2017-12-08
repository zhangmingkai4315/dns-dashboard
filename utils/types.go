package utils

// Global config part
type Global struct {
	Type   string `ini:"type"`
	Server string `ini:"server"`
	Port   int    `ini:"port"`
	Source string `ini:"source"`
	Path   string `ini:"path"`
	Grok   string `ini:"grok"`
}

// DNS config part
type DNS struct {
	Grok string `ini:"grok"`
}

// Config from config.ini
type Config struct {
	Global
	DNS
}
