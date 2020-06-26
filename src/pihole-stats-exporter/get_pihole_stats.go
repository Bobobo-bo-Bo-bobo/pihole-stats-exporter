package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func getPiHoleRawSummary(request *http.Request) (PiHoleRawSummary, error) {
	var rawsum PiHoleRawSummary

	// get raw summary
	result, err := fetchPiHoleData(config, "summaryRaw")
	if err != nil {
		log.WithFields(log.Fields{
			"remote_address": request.RemoteAddr,
			"error":          err.Error(),
			"pihole_request": "summaryRaw",
		}).Error(formatLogString("Can't fetch data from PiHole server"))

		return rawsum, err
	}

	if result.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"remote_address": request.RemoteAddr,
			"status_code":    result.StatusCode,
			"status":         result.Status,
			"pihole_request": "summaryRaw",
		}).Error(formatLogString("Unexpected HTTP status from PiHole server"))

		return rawsum, fmt.Errorf("Unexpected HTTP status from PiHole server")
	}

	err = json.Unmarshal(result.Content, &rawsum)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"pihole_request": "summaryRaw",
		}).Error(formatLogString("Can't decode received result as JSON data"))

		return rawsum, err
	}

	return rawsum, nil
}

func getPiHoleQueryTypes(request *http.Request) (PiHoleQueryTypes, error) {
	var qtypes PiHoleQueryTypes

	// get DNS queries by type
	result, err := fetchPiHoleData(config, "getQueryTypes")
	if err != nil {
		log.WithFields(log.Fields{
			"remote_address": request.RemoteAddr,
			"error":          err.Error(),
			"pihole_request": "getQueryTypes",
		}).Error(formatLogString("Can't fetch data from PiHole server"))

		return qtypes, err
	}

	if result.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"remote_address": request.RemoteAddr,
			"status_code":    result.StatusCode,
			"status":         result.Status,
			"pihole_request": "getQueryTypes",
		}).Error(formatLogString("Unexpected HTTP status from PiHole server"))

		return qtypes, fmt.Errorf("Unexpected HTTP status from PiHole server")
	}

	err = json.Unmarshal(result.Content, &qtypes)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"pihole_request": "getQueryTypes",
		}).Error(formatLogString("Can't decode received result as JSON data"))

		return qtypes, err
	}

	return qtypes, nil
}
