package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
	"github.com/zhangmingkai4315/dns-dashboard/web"
)

var (
	configFile = flag.String("c", "", "config file path")
)
var usage = `Usage: dns-dashboard [options...] 
Options:
  -c       config file path
`

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()
	if *configFile == "" {
		usageAndExit("config file must provided")
	}
	config, err := utils.LoadConfigFromFile("./config.ini")
	if err != nil {
		usageAndExit(fmt.Sprintf("load config file err:%s", err))
	}
	if config.Type == "master" {
		if config.Server != "" && config.Port > 0 {
			web.StartServer(config)
		} else {
			usageAndExit("server and port config error")
		}
	}

}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, "Error: %s", msg)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	os.Exit(1)
}
