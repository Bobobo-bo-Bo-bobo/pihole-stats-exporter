package main

import (
	"net/http"
	"time"
)

// PiHoleRawSummary - raw summary
type PiHoleRawSummary struct {
	DomainsBeingBlocked uint64                   `json:"domains_being_blocked"`
	DNSQueriesToday     uint64                   `json:"dns_queries_today"`
	AdsBlockedToday     uint64                   `json:"ads_blocked_today"`
	AdsPercentageToday  uint64                   `json:"ads_percentage_today"`
	UniqueDomains       uint64                   `json:"unique_domains"`
	QueriesForwarded    uint64                   `json:"queries_forwarded"`
	QueriesCached       uint64                   `json:"queries_cached"`
	ClientsEverSeend    uint64                   `json:"clients_ever_seen"`
	UniqueClients       uint64                   `json:"unique_clients"`
	DNSQueriesAllTypes  uint64                   `json:"dns_queries_all_types"`
	ReplyNODATA         uint64                   `json:"reply_NODATA"`
	ReplyNXDOMAIN       uint64                   `json:"reply_NXDOMAIN"`
	ReplyCNAME          uint64                   `json:"reply_CNAME"`
	ReplyIP             uint64                   `json:"reply_IP"`
	PrivacyLevel        uint                     `json:"privacy_level"`
	Status              string                   `json:"status"`
	GravityLastUpdated  PiHoleGravityLastUpdated `json:"gravity_last_updated"`
}

// PiHoleGravityLastUpdated - information about last gravity update
type PiHoleGravityLastUpdated struct {
	FileExists bool                             `json:"file_exists"`
	Absolute   uint64                           `json:"absolute"`
	Relative   PiHoleGravityLastUpdatedRelative `json:"relative"`
}

// PiHoleGravityLastUpdatedRelative - relative time of last gravity update
type PiHoleGravityLastUpdatedRelative struct {
	Days    uint64 `json:"days"`
	Hours   uint64 `json:"hours"`
	Minutes uint64 `json:"minutes"`
}

// Configuration - hold configuration information
type Configuration struct {
	PiHole   PiHoleConfiguration
	Exporter ExporterConfiguration
}

// PiHoleConfiguration - Configure access to PiHole
type PiHoleConfiguration struct {
	URL            string `ini:"url"`
	AuthHash       string `ini:"auth"`
	InsecureSSL    bool   `ini:"insecure_ssl"`
	CAFile         string `ini:"ca_file"`
	Timeout        uint   `ini:"timeout"`
	FollowRedirect bool   `ini:"follow_redirect"`
	timeout        time.Duration
}

// ExporterConfiguration - configure metric exporter
type ExporterConfiguration struct {
	URL            string `ini:"url"`
	PrometheusPath string `ini:"prometheus_path"`
	InfluxDataPath string `ini:"influxdata_path"`
	SSLCert        string `ini:"ssl_cert"`
	SSLKey         string `ini:"ssl_key"`
}

// HTTPResult - result of the http_request calls
type HTTPResult struct {
	URL        string
	StatusCode int
	Status     string
	Header     http.Header
	Content    []byte
}
