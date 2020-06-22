package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	var help = flag.Bool("help", false, "Show help text")
	var version = flag.Bool("version", false, "Show version information")
	var configFile = flag.String("config", "", "Path to configuration file")
	var logFmt = new(log.TextFormatter)

	flag.Usage = showUsage

	flag.Parse()

	logFmt.FullTimestamp = true
	logFmt.TimestampFormat = time.RFC3339
	log.SetFormatter(logFmt)

	if len(flag.Args()) > 0 {
		log.Fatal(formatLogString("Trailing arguments"))
	}

	if *help {
		showUsage()
		os.Exit(0)
	}

	if *version {
		showVersion()
		os.Exit(0)
	}

	if *configFile == "" {
		log.Fatal(formatLogString("Path to configuration file (--config) is mandatory"))
	}

	config, err := parseConfigurationFile(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"config_file": *configFile,
			"error":       err.Error(),
		}).Fatal(formatLogString("Can't parse configuration file"))
	}

	if config.Exporter.PrometheusPath == "" && config.Exporter.InfluxDataPath == "" {
		log.WithFields(log.Fields{
			"config_file": *configFile,
		}).Fatal(formatLogString("Neither the path for Prometheus metrics nor for InfluxDB metrics are set"))
	}

	if config.Exporter.PrometheusPath == "" {
		log.WithFields(log.Fields{
			"config_file":     *configFile,
			"prometheus_path": config.Exporter.PrometheusPath,
			"influxdata_path": config.Exporter.InfluxDataPath,
		}).Warning(formatLogString("Path for Prometheus metrics is not set, disabling Prometheus metrics"))
	}

	if config.Exporter.InfluxDataPath == "" {
		log.WithFields(log.Fields{
			"config_file":     *configFile,
			"prometheus_path": config.Exporter.PrometheusPath,
			"influxdata_path": config.Exporter.InfluxDataPath,
		}).Warning(formatLogString("Path for InfluxDB metrics is not set, disabling InfluxDB metrics"))
	}

	// spawn HTTP server
	_uri, err := url.Parse(config.Exporter.URL)
	if err != nil {
		log.WithFields(log.Fields{
			"config_file":  *configFile,
			"exporter_url": config.Exporter.URL,
			"error":        err.Error(),
		}).Fatal(formatLogString("Can't parse exporter URL"))
	}

	// XXX: This should go into validateConfiguration
	if _uri.Scheme != "http" && _uri.Scheme != "https" {
		log.WithFields(log.Fields{
			"config_file":  *configFile,
			"exporter_url": config.Exporter.URL,
		}).Fatal(formatLogString("Invalid or unsupported URL scheme"))
	}

	router := mux.NewRouter()
	subRouterGet := router.Methods("GET").Subrouter()

	if config.Exporter.PrometheusPath != "" {
		subRouterGet.HandleFunc(config.Exporter.PrometheusPath, prometheusExporter)
	}

	if config.Exporter.InfluxDataPath != "" {
		subRouterGet.HandleFunc(config.Exporter.InfluxDataPath, influxExporter)
	}

	log.WithFields(log.Fields{
		"config_file":     *configFile,
		"exporter_url":    config.Exporter.URL,
		"prometheus_path": config.Exporter.PrometheusPath,
		"influxdata_path": config.Exporter.InfluxDataPath,
	}).Info(formatLogString("Starting HTTP listener"))

	router.Host(_uri.Host)

	// XXX: Add timeout values to configuration file instead of hardcoding values
	httpSrv := &http.Server{
		Addr:         _uri.Host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	if _uri.Scheme == "https" {
		_tls, err := generateTLSConfiguration(config)
		if err != nil {
			log.WithFields(log.Fields{
				"config_file":  *configFile,
				"exporter_url": config.Exporter.URL,
				"error":        err.Error(),
			}).Fatal(formatLogString("Can't create TLS context"))
		}

		httpSrv.TLSConfig = _tls
		httpSrv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
	}

	// start HTTP routine and wait for termination signals to arrive
	go func() {
		if _uri.Scheme == "https" {
			err = httpSrv.ListenAndServeTLS(config.Exporter.SSLCert, config.Exporter.SSLKey)
		} else {
			err = httpSrv.ListenAndServe()
		}
		if err != nil {
			log.WithFields(log.Fields{
				"config_file":  *configFile,
				"exporter_url": config.Exporter.URL,
				"error":        err.Error(),
			}).Fatal(formatLogString("Can't start HTTP server"))
		}
	}()

	// listen for signals
	sigChan := make(chan os.Signal, 1)

	// Listen for SIGINT, SIGKILL and SIGTERM signals
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	sig := <-sigChan

	log.WithFields(log.Fields{
		"config_file":     *configFile,
		"exporter_url":    config.Exporter.URL,
		"prometheus_path": config.Exporter.PrometheusPath,
		"influxdata_path": config.Exporter.InfluxDataPath,
		"signal":          sig.String(),
	}).Info(formatLogString("Received termination signal, terminating HTTP server"))

	_ctx, _cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer _cancel()

	// This will shutdown the server immediately if no connection is present, otherwise wait for 15 seconds
	httpSrv.Shutdown(_ctx)

	os.Exit(0)
}
