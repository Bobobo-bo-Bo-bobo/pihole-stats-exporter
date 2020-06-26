package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func influxExporter(response http.ResponseWriter, request *http.Request) {
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
	now := time.Now().Unix() * 1e+09

	payload = []byte(fmt.Sprintf(`pihole,type=summary,upstream=%s,type=domains_being_blocked value=%d %d
pihole,type=summary,upstream=%s,type=dns_queries_today value=%d %d
pihole,type=summary,upstream=%s,type=ads_blocked_today value=%d %d
pihole,type=summary,upstream=%s,type=ads_percentage_today value=%f %d
pihole,type=summary,upstream=%s,type=unique_domains value=%d %d
pihole,type=summary,upstream=%s,type=queries_forwarded value=%d %d
pihole,type=summary,upstream=%s,type=queries_cached value=%d %d
pihole,type=summary,upstream=%s,type=clients_ever_seen value=%d %d
pihole,type=summary,upstream=%s,type=unique_clients value=%d %d
pihole,type=summary,upstream=%s,type=dns_queries_all_types value=%d %d
pihole,type=summary,upstream=%s,type=reply_NODATA value=%d %d
pihole,type=summary,upstream=%s,type=reply_NXDOMAIN value=%d %d
pihole,type=summary,upstream=%s,type=reply_CNAME value=%d %d
pihole,type=summary,upstream=%s,type=reply_IP value=%d %d
pihole,type=summary,upstream=%s,type=privacy_level value=%d %d
pihole,type=querytypes,upstream=%s,type=A value=%f %d
pihole,type=querytypes,upstream=%s,type=AAAA value=%f %d
pihole,type=querytypes,upstream=%s,type=ANY value=%f %d
pihole,type=querytypes,upstream=%s,type=SRV value=%f %d
pihole,type=querytypes,upstream=%s,type=SOA value=%f %d
pihole,type=querytypes,upstream=%s,type=PTR value=%f %d
pihole,type=querytypes,upstream=%s,type=TXT value=%f %d
pihole,type=querytypes,upstream=%s,type=NAPTR value=%f %d
`,
		config.PiHole.URL, rawsum.DomainsBeingBlocked, now,
		config.PiHole.URL, rawsum.DNSQueriesToday, now,
		config.PiHole.URL, rawsum.AdsBlockedToday, now,
		config.PiHole.URL, rawsum.AdsPercentageToday, now,
		config.PiHole.URL, rawsum.UniqueDomains, now,
		config.PiHole.URL, rawsum.QueriesForwarded, now,
		config.PiHole.URL, rawsum.QueriesCached, now,
		config.PiHole.URL, rawsum.ClientsEverSeend, now,
		config.PiHole.URL, rawsum.UniqueClients, now,
		config.PiHole.URL, rawsum.DNSQueriesAllTypes, now,
		config.PiHole.URL, rawsum.ReplyNODATA, now,
		config.PiHole.URL, rawsum.ReplyNXDOMAIN, now,
		config.PiHole.URL, rawsum.ReplyCNAME, now,
		config.PiHole.URL, rawsum.ReplyIP, now,
		config.PiHole.URL, rawsum.PrivacyLevel, now,
		config.PiHole.URL, qtypes.Querytypes.A, now,
		config.PiHole.URL, qtypes.Querytypes.AAAA, now,
		config.PiHole.URL, qtypes.Querytypes.ANY, now,
		config.PiHole.URL, qtypes.Querytypes.SRV, now,
		config.PiHole.URL, qtypes.Querytypes.SOA, now,
		config.PiHole.URL, qtypes.Querytypes.PTR, now,
		config.PiHole.URL, qtypes.Querytypes.TXT, now,
		config.PiHole.URL, qtypes.Querytypes.NAPTR, now,
	))
	response.Write(payload)

	// discard slice and force gc to free the allocated memory
	payload = nil
}
