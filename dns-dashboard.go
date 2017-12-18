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
	// 如果配置文件未设置,退出程序
	if *configFile == "" {
		utils.UsageAndExit("config file must provide")
	}
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		utils.UsageAndExit("config file not exist")
	}
	config, err := utils.LoadConfigFromFile(*configFile)
	if err != nil {
		utils.UsageAndExit(fmt.Sprintf("load config file err:%s", err))
	} else {
		log.Println("load config file sucess")
	}
	//日志文件可以从配置文件中获取
	logFile := *queryLogFile
	if *queryLogFile == "" {
		logFile = config.DNS.Path
	}

	if logFile == "" {
		utils.UsageAndExit("log file must provide")
	}

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		utils.UsageAndExit("log file not exist")
	}

	manager, err := analyzer.NewDNSStatsManager(logFile, config.DNS.Grok)
	if err != nil {
		utils.UsageAndExit(fmt.Sprintf("can't create manager:%s", err))
	}

	db, err := model.GetDB()
	if err != nil {
		utils.UsageAndExit("DB instance not ready")
	}
	// 启动anylyzer每隔五秒将结果写入数据库
	go func() {
		for {
			log.Println("Start the dns analyzer processing...")
			manager.Start()
			time.Sleep(time.Second * 5)
			log.Println("Stop the processing...")
			stats, err := manager.Stop()
			if err != nil {
				log.Errorf("Error when stop the manager %s", err)
				continue
			}
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
			// 将处理后的数据存入数据库，由于无法直接存入数组，使用string类型保存统计数据
			db.Create(&serial)
		}
	}()
	// 暂时只实现了master的功能，未来计划通过agent的方式将结果推送到主控，
	// 实现数据中心的集中化管理
	if config.Type == "master" {
		if config.Server != "" && config.Port > 0 {
			web.StartServer(config)
		} else {
			utils.UsageAndExit("server and port config error")
		}
	}
}
