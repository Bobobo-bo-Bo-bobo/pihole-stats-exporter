package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func prometheusExporter(response http.ResponseWriter, request *http.Request) {
	var rawsum PiHoleRawSummary
	var qtypes PiHoleQueryTypes
	var err error
	var payload []byte

	log.WithFields(log.Fields{
		"method":         request.Method,
		"url":            request.URL.String(),
		"protocol":       request.Proto,
		"host":           request.Host,
		"remote_address": request.RemoteAddr,
		"headers":        fmt.Sprintf("%+v\n", request.Header),
	}).Info(formatLogString("HTTP request from client received"))

	response.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")

	// get raw summary
	rawsum, err = getPiHoleRawSummary(request)
	if err != nil {
		log.WithFields(log.Fields{
			"remote_address": request.RemoteAddr,
			"error":          err.Error(),
			"pihole_request": "summaryRaw",
		}).Error(formatLogString("Can't fetch data from PiHole server"))

		response.Write([]byte("502 bad gateway"))
		response.WriteHeader(http.StatusBadGateway)

		return
	}

	// get DNS queriey by type
	qtypes, err = getPiHoleQueryTypes(request)
	if err != nil {
		log.WithFields(log.Fields{
			"remote_address": request.RemoteAddr,
			"error":          err.Error(),
			"pihole_request": "getQueryTypes",
		}).Error(formatLogString("Can't fetch data from PiHole server"))

		response.Write([]byte("502 bad gateway"))
		response.WriteHeader(http.StatusBadGateway)

		return
	}

	payload = []byte(fmt.Sprintf(`#HELP pihole_domains_blocked_total Number of blocked domains
#TYPE pihole_domains_blocked_total counter
pihole_domains_blocked_total{upstream="%s"} %d
#HELP pihole_dns_queries_today_total Number of DNS queries received today
#TYPE pihole_dns_queries_today_total counter
pihole_dns_queries_today_total{upstream="%s"} %d
#HELP pihole_ads_today_total Number if requests blackholed
#TYPE pihole_ads_today_total counter
pihole_ads_today_total{upstream="%s"} %d
#HELP pihole_ads_today_ratio Percentage of blackholed requests
#TYPE pihole_ads_today_ratio gauge
pihole_ads_today_ratio{upstream="%s"} %f
#HELP pihole_unique_domains_total Unique domains seen today
#TYPE pihole_unique_domains_total counter
pihole_unique_domains_total{upstream="%s"} %d
#HELP pihole_queries_forwarded Number of DNS requests forwarded
#TYPE pihole_queries_forwarded gauge
pihole_queries_forwarded{upstream="%s"} %d
#HELP pihole_queries_cached Number of DNS requests cached
#TYPE pihole_queries_cached gauge
pihole_queries_cached{upstream="%s"} %d
#HELP pihole_clients_ever_seen_total Number of clients ever seen
#TYPE pihole_clients_ever_seen_total counter
pihole_clients_ever_seen_total{upstream="%s"} %d
#HELP pihole_unique_clients Number of unique clients
#TYPE pihole_unique_clients gauge
pihole_unique_clients{upstream="%s"} %d
#HELP pihole_dns_queries_all_types_total Number of DNS queries of all types
#TYPE pihole_dns_queries_all_types_total counter
pihole_dns_queries_all_types_total{upstream="%s"} %d
#HELP pihole_reply_total DNS replies by type
#TYPE pihole_reply_total counter
pihole_reply_total{upstream="%s",reply="NODATA"} %d
pihole_reply_total{upstream="%s",reply="NXDOMAIN"} %d
pihole_reply_total{upstream="%s",reply="CNAME"} %d
pihole_reply_total{upstream="%s",reply="IP"} %d
#HELP pihole_privacy_level PiHole privacy level
#TYPE pihole_privacy_level gauge
pihole_privacy_level{upstream="%s"} %d
#HELP pihole_query_type_ratio Ratio of DNS type requested from clients
#TYPE pihole_query_type_ratio gauge
pihole_query_type_ratio{upstream="%s",type="A"} %f
pihole_query_type_ratio{upstream="%s",type="AAAA"} %f
pihole_query_type_ratio{upstream="%s",type="ANY"} %f
pihole_query_type_ratio{upstream="%s",type="SRV"} %f
pihole_query_type_ratio{upstream="%s",type="SOA"} %f
pihole_query_type_ratio{upstream="%s",type="PTR"} %f
pihole_query_type_ratio{upstream="%s",type="TXT"} %f
pihole_query_type_ratio{upstream="%s",type="NAPTR"} %f
`,
		config.PiHole.URL, rawsum.DomainsBeingBlocked,
		config.PiHole.URL, rawsum.DNSQueriesToday,
		config.PiHole.URL, rawsum.AdsBlockedToday,
		config.PiHole.URL, rawsum.AdsPercentageToday/100.0,
		config.PiHole.URL, rawsum.UniqueDomains,
		config.PiHole.URL, rawsum.QueriesForwarded,
		config.PiHole.URL, rawsum.QueriesCached,
		config.PiHole.URL, rawsum.ClientsEverSeend,
		config.PiHole.URL, rawsum.UniqueClients,
		config.PiHole.URL, rawsum.DNSQueriesAllTypes,
		config.PiHole.URL, rawsum.ReplyNODATA,
		config.PiHole.URL, rawsum.ReplyNXDOMAIN,
		config.PiHole.URL, rawsum.ReplyCNAME,
		config.PiHole.URL, rawsum.ReplyIP,
		config.PiHole.URL, rawsum.PrivacyLevel,
		config.PiHole.URL, float64(qtypes.Querytypes.A)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.AAAA)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.ANY)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.SRV)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.SOA)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.PTR)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.TXT)/100.0,
		config.PiHole.URL, float64(qtypes.Querytypes.NAPTR)/100.0,
	))
	response.Write(payload)

	// discard slice and force gc to free the allocated memory
	payload = nil
}
