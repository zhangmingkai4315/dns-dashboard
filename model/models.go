package model

// DNSSerialData will save data to database
type DNSSerialData struct {
	ID             uint    `json:"id"`
	Timestamp      string  `json:"timestamp"`
	DomainStats    string  `json:"domain_stats"`
	SubDomainStats string  `json:"sub_domain_stats"`
	Duration       float64 `json:"duration"`
	TypeStats      string  `json:"type_stats"`
	IPStats        string  `json:"ip_stats"`
}
