package main

import (
	"flag"
	"os"
	"time"

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

	os.Exit(0)
}
