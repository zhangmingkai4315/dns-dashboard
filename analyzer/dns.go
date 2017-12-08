package analyzer

import (
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/vjeantet/grok"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
)

var (
	// MaxWorkerNumber define the max size of worker number
	MaxWorkerNumber = 40
	// TopCounts define the top 10 results
	TopCounts = 10
)

var g *grok.Grok

func init() {
	g, _ = grok.New()
}

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
	IPStats        []IPInfo     `json:"ip_stats"`
	DomainStats    []DomainInfo `json:"domain_stats"`
	SubDomainStats []DomainInfo `json:"sub_domain_stats"`
	// process duration
	Duration  float64    `json:"duration"`
	TypeStats []TypeInfo `json:"type_stats"`
}

// DNSStatsManager define the manager to
type DNSStatsManager struct {
	mutex       sync.Mutex
	file        *tail.Tail
	rawChannel  chan string
	stopChannel chan struct{}
	startTime   time.Time
	Timestamp   time.Time
	worker      int
	processing  bool
	stats       DNSStats
	data        []map[string]map[string]int
	grok        string
}

// NewDNSStatsManager function will generate a new manager
// to manage the analyz process
func NewDNSStatsManager(file string, grok string) (*DNSStatsManager, error) {
	location := &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	tailFile, err := tail.TailFile(file, tail.Config{Follow: true, Location: location})
	if err != nil {
		return nil, err
	}
	rawChannel := make(chan string, 1000)
	stopChannel := make(chan struct{})
	var data []map[string]map[string]int
	for i := 0; i < MaxWorkerNumber; i++ {
		data = append(data, make(map[string]map[string]int))
	}
	if err != nil {
		return nil, err
	}
	return &DNSStatsManager{
		file:        tailFile,
		rawChannel:  rawChannel,
		processing:  false,
		stopChannel: stopChannel,
		worker:      MaxWorkerNumber,
		data:        data,
		grok:        grok,
	}, nil
}

// Start will start to get data from file
// and start all works
func (manager *DNSStatsManager) Start() {
	manager.data = make([]map[string]map[string]int, MaxWorkerNumber)

	for i := 0; i < manager.worker; i++ {
		go func(index int) {
			manager.mutex.Lock()
			defer manager.mutex.Unlock()
			manager.data[index] = make(map[string]map[string]int)
			manager.data[index]["domain"] = make(map[string]int)
			manager.data[index]["ip"] = make(map[string]int)
			manager.data[index]["type"] = make(map[string]int)
			manager.data[index]["subdomain"] = make(map[string]int)
			domainMap := manager.data[index]["domain"]
			subDomainMap := manager.data[index]["subdomain"]
			ipMap := manager.data[index]["ip"]
			typeMap := manager.data[index]["type"]
			// process worker will do all data process
			for {
				select {
				case rawText := <-manager.rawChannel:
					// regex process
					rawInfo, err := manager.getRawFromText(rawText)
					if err != nil {
						log.Println(err)
						continue
					}
					domainMap[rawInfo.Domain] = domainMap[rawInfo.Domain] + 1
					subdomain := utils.RemoveSubDomain(rawInfo.Domain)
					subDomainMap[subdomain] = subDomainMap[subdomain] + 1
					ipMap[rawInfo.IP] = ipMap[rawInfo.IP] + 1
					typeMap[rawInfo.Type] = typeMap[rawInfo.Type] + 1
				case <-manager.stopChannel:
					return
				}
			}
		}(i)
	}

	manager.startTime = time.Now()
	for line := range manager.file.Lines {
		manager.rawChannel <- line.Text
		log.Println(line.Text)
		select {
		case <-manager.stopChannel:
			log.Println("reader recive stop signal")
			return
		default:
			continue
		}
	}
}

// Stop will send stop signal to manage and stop process
func (manager *DNSStatsManager) Stop() (*DNSStats, error) {
	// close read and process worker
	log.Println("Sending stop signal to all worker and reader goroutine")
	for i := 0; i < manager.worker+1; i++ {
		manager.stopChannel <- struct{}{}
	}
	duration := time.Now().Sub(manager.startTime).Seconds()
	manager.Timestamp = time.Now()
	// counts top results
	domainCounts := make(map[string]int)
	ipCounts := make(map[string]int)
	typeCounts := make(map[string]int)
	subDomainCounts := make(map[string]int)
	for i := 0; i < manager.worker; i++ {
		for k, v := range manager.data[i]["domain"] {
			domainCounts[k] = domainCounts[k] + v
		}
		for k, v := range manager.data[i]["ip"] {
			ipCounts[k] = ipCounts[k] + v
		}
		for k, v := range manager.data[i]["type"] {
			typeCounts[k] = typeCounts[k] + v
		}
		for k, v := range manager.data[i]["subdomain"] {
			subDomainCounts[k] = subDomainCounts[k] + v
		}
	}

	topDomain := make([]DomainInfo, 0)
	for k, v := range domainCounts {
		topDomain = append(topDomain, DomainInfo{k, v})
	}
	sort.Slice(topDomain, func(i, j int) bool {
		return topDomain[i].Sum > topDomain[j].Sum
	})
	if len(topDomain) > TopCounts {
		topDomain = topDomain[:TopCounts]
	}
	var topIP []IPInfo
	for k, v := range ipCounts {
		topIP = append(topIP, IPInfo{k, v})
	}
	sort.Slice(topIP, func(i, j int) bool {
		return topIP[i].Sum > topIP[j].Sum
	})
	if len(topIP) > TopCounts {
		topIP = topIP[:TopCounts]
	}

	var topType []TypeInfo
	for k, v := range typeCounts {
		topType = append(topType, TypeInfo{k, v})
	}
	sort.Slice(topType, func(i, j int) bool {
		return topType[i].Sum > topType[j].Sum
	})
	if len(topType) > TopCounts {
		topType = topType[:TopCounts]
	}

	var topSubDomain []DomainInfo
	for k, v := range subDomainCounts {
		topSubDomain = append(topSubDomain, DomainInfo{k, v})
	}
	sort.Slice(topSubDomain, func(i, j int) bool {
		return topSubDomain[i].Sum > topSubDomain[j].Sum
	})
	if len(topSubDomain) > TopCounts {
		topSubDomain = topSubDomain[:TopCounts]
	}
	stats := &DNSStats{
		Duration:       duration,
		DomainStats:    topDomain,
		SubDomainStats: topSubDomain,
		IPStats:        topIP,
		TypeStats:      topType,
	}
	return stats, nil
}

func (manager *DNSStatsManager) getRawFromText(text string) (*RawInfo, error) {
	values, err := g.Parse(manager.grok, text)
	if err != nil {
		return nil, err
	}
	rawInfo := &RawInfo{
		Domain: values["domain"],
		IP:     values["ip"],
		Type:   values["type"],
	}
	return rawInfo, nil
}
