package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zhangmingkai4315/dns-dashboard/model"

	log "github.com/sirupsen/logrus"
	"github.com/zhangmingkai4315/dns-dashboard/analyzer"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
	"github.com/zhangmingkai4315/dns-dashboard/web"
)

var (
	configFile   = flag.String("c", "", "config file path")
	queryLogFile = flag.String("f", "", "dns query log file path")
)
var usage = `Usage: dns-dashboard [options...] 
Options:
  -c       config file path
  -f       file for read dns log stream
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()
	if *configFile == "" {
		utils.UsageAndExit("config file must provided")
	}
	config, err := utils.LoadConfigFromFile(*configFile)
	if err != nil {
		utils.UsageAndExit(fmt.Sprintf("load config file err:%s", err))
	}

	if *queryLogFile == "" {
		utils.UsageAndExit("query file must provided")
	}
	manager, err := analyzer.NewDNSStatsManager(*queryLogFile, config.DNS.Grok)
	if err != nil {
		utils.UsageAndExit(fmt.Sprintf("can't create manager:%s", err))
	}

	ticker := time.NewTicker(time.Millisecond)
	db, err := model.GetDB()
	if err != nil {
		utils.UsageAndExit("DB instance not ready")
	}
	go func() {
		for _ = range ticker.C {
			log.Println("Start the dns analyzer processing...")
			go manager.Start()
			time.Sleep(time.Second * 5)
			log.Println("Stop the processing")
			stats, err := manager.Stop()
			if err == nil {
				var domainStats string
				domainStatsBytes, err := json.Marshal(stats.DomainStats)
				if err == nil {
					domainStats = string(domainStatsBytes)
				}
				var ipStats string
				ipStatsBytes, err := json.Marshal(stats.IPStats)
				if err == nil {
					ipStats = string(ipStatsBytes)
				}

				var subDomainStats string
				subDomainStatsBytes, err := json.Marshal(stats.SubDomainStats)
				if err == nil {
					subDomainStats = string(subDomainStatsBytes)
				}
				var typeStats string
				typeStatsBytes, err := json.Marshal(stats.TypeStats)
				if err == nil {
					typeStats = string(typeStatsBytes)
				}
				serial := model.DNSSerialData{
					Timestamp:      manager.Timestamp.String(),
					DomainStats:    domainStats,
					IPStats:        ipStats,
					Duration:       stats.Duration,
					SubDomainStats: subDomainStats,
					TypeStats:      typeStats,
				}
				db.Create(&serial)
			}
		}

	}()

	if config.Type == "master" {
		if config.Server != "" && config.Port > 0 {
			web.StartServer(config)
		} else {
			utils.UsageAndExit("server and port config error")
		}
	}
}
