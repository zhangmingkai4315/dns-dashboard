package analyzer

import (
	"time"

	"github.com/hpcloud/tail"
	"github.com/vjeantet/grok"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
)

var (
	// MaxBufferLine define the max size of read lines
	MaxBufferLine = 500
	// MaxWorker define the max size of worker number
	MaxWorkerNumber = 10
)

// IPInfo define the statics data based ip infomation
type IPInfo struct {
	IP  string `json:"ip"`
	Sum int    `json:"sum"`
}

// DomainInfo define the statics data based domain infomation
type DomainInfo struct {
	Domain string `json:"domain"`
	Sum    int    `json:"sum"`
}

// TypeInfo define the statics data based type infomation
type TypeInfo struct {
	Type string `json:"type"`
	Sum  int    `json:"sum"`
}

// RawInfo Using regex to get from dns log
type RawInfo struct {
	Domain string
	Type   string
	IP     string
}

// DNSStats define the total stats of everytimes
// generater by core analyzer
type DNSStats struct {
	IPStats     []IPInfo     `json:"ip_stats"`
	DomainStats []DomainInfo `json:"domain_stats"`
	TimeStamp   string       `json:"timestamp"`
	TypeStats   []TypeInfo   `json:"type_stats"`
}

// DNSStatsManager define the manager to
type DNSStatsManager struct {
	file          *tail.Tail
	lastTimeStamp time.Time
	rawChannel    chan string
	stopChannel   chan struct{}
	worker        int
	processing    bool
	stats         DNSStats
	data          []map[string]map[string]int
	g             *grok.Grok
	config        *utils.Config
}

// NewDNSStatsManager function will generate a new manager
// to manage the analyz process
func NewDNSStatsManager(file string) (*DNSStatsManager, error) {
	tailFile, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}
	rawChannel := make(chan string)
	stopChannel := make(chan struct{})
	data := make([]map[string]map[string]int, MaxWorkerNumber)
	g, err := grok.New()
	if err != nil {
		return nil, err
	}
	config, err := utils.GetConfig()
	if err != nil {
		return nil, err
	}
	return &DNSStatsManager{
		file:          tailFile,
		lastTimeStamp: time.Now(),
		rawChannel:    rawChannel,
		processing:    false,
		stopChannel:   stopChannel,
		worker:        MaxWorkerNumber,
		data:          data,
		g:             g,
		config:        config,
	}, nil
}

// Start will start to get data from file
// and start all works
func (manager *DNSStatsManager) Start() {
	for i := 0; i < manager.worker; i++ {
		go func() {
			domainMap := manager.data[i]["domain"]
			ipMap := manager.data[i]["ip"]
			typeMap := manager.data[i]["type"]
			// process worker will do all data process
			for {
				select {
				case rawText := <-manager.rawChannel:
					// regex process
					rawInfo, err := manager.getRawFromText(rawText)
					if err != nil {
						continue
					}
					domainMap[rawInfo.Domain] = domainMap[rawInfo.Domain] + 1
					ipMap[rawInfo.IP] = ipMap[rawInfo.IP] + 1
					typeMap[rawInfo.Type] = typeMap[rawInfo.Type] + 1
				}
			}
		}()
	}
	for line := range manager.file.Lines {
		manager.rawChannel <- line.Text
		select {
		case <-manager.stopChannel:
			break
		}
	}
}

// Stop will send stop signal to manage and stop process
func (manager *DNSStatsManager) Stop() {
	manager.stopChannel <- struct{}{}
}

func (manager *DNSStatsManager) getRawFromText(text string) (*RawInfo, error) {
	grok := manager.config.DNS.Grok
	values, err := manager.g.Parse(grok, text)
	if err != nil {
		return nil, err
	}
	rawInfo := &RawInfo{
		Domain: values["domain"],
		IP:     values["client"],
		Type:   values["type"],
	}
	return rawInfo, nil
}
