package main

import (
	"crypto/tls"
	"fmt"
)

func generateTLSConfiguration(cfg *Configuration) (*tls.Config, error) {
	if cfg.Exporter.SSLCert == "" || cfg.Exporter.SSLKey == "" {
		return nil, fmt.Errorf("No server certificate files provided")
	}

	result := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	return result, nil
}
