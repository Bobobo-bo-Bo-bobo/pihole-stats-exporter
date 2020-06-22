package main

import (
	"fmt"

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

	return &config, nil
}

func validateConfiguration(cfg Configuration) error {
	if cfg.PiHole.URL == "" {
		return fmt.Errorf("URL to PiHole is missing")
	}

	return nil
}
