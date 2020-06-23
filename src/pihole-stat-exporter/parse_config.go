package main

import (
	"fmt"
	"time"

	ini "gopkg.in/ini.v1"
)

func parseConfigurationFile(f string) (*Configuration, error) {
	var err error

	config := Configuration{
		Exporter: ExporterConfiguration{
			URL:            defaultExporterURL,
			PrometheusPath: defaultPrometheusPath,
			InfluxDataPath: defaultInfluxDataPath,
		},
		PiHole: PiHoleConfiguration{
			Timeout: 15,
		},
	}

	cfg, err := ini.Load(f)
	if err != nil {
		return nil, err
	}

	pihole, err := cfg.GetSection("pihole")
	if err != nil {
		return nil, err
	}
	err = pihole.MapTo(&config.PiHole)
	if err != nil {
		return nil, err
	}

	exporter, err := cfg.GetSection("exporter")
	if err != nil {
		return nil, err
	}
	err = exporter.MapTo(&config.Exporter)
	if err != nil {
		return nil, err
	}

	err = validateConfiguration(config)
	if err != nil {
		return nil, err
	}

	config.PiHole.timeout = time.Duration(config.PiHole.Timeout) * time.Second

	return &config, nil
}

func validateConfiguration(cfg Configuration) error {
	if cfg.PiHole.URL == "" {
		return fmt.Errorf("URL to PiHole is missing")
	}
	if cfg.PiHole.Timeout == 0 {
		return fmt.Errorf("Invalid timeout")
	}

	if cfg.Exporter.PrometheusPath != "" && cfg.Exporter.PrometheusPath[0] != '/' {
		return fmt.Errorf("Prometheus path must be an absolute path")
	}

	if cfg.Exporter.InfluxDataPath != "" && cfg.Exporter.InfluxDataPath[0] != '/' {
		return fmt.Errorf("InfluxDB path must be an absolute path")
	}
	return nil
}
