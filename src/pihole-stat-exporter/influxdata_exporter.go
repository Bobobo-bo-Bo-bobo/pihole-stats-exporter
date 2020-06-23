package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func influxExporter(response http.ResponseWriter, request *http.Request) {
	log.WithFields(log.Fields{
		"method":         request.Method,
		"url":            request.URL.String(),
		"protocol":       request.Proto,
		"host":           request.Host,
		"remote_address": request.RemoteAddr,
		"headers":        fmt.Sprintf("%+v\n", request.Header),
	}).Info(formatLogString("HTTP request from client received"))

	// get raw summary
	result, err := fetchPiHoleData(config, "summaryRaw")
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"pihole_request": "summaryRaw",
		}).Error(formatLogString("Can't fetch data from PiHole server"))

		result.Header.Add("X-Clacks-Overhead", "GNU Terry Pratchett")

		if result.StatusCode != 0 {
			response.WriteHeader(result.StatusCode)
		} else {
			response.WriteHeader(http.StatusBadGateway)
		}

		return
	}

}
