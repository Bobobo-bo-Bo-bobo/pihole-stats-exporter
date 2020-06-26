package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func fetchPiHoleData(cfg *Configuration, stat string) (HTTPResult, error) {
	var result HTTPResult
	var transp *http.Transport
	var err error
	var piurl string

	_url, err := url.Parse(cfg.PiHole.URL)
	if err != nil {
		return result, err
	}

	if _url.Scheme == "https" {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
		if cfg.PiHole.InsecureSSL {
			transp.TLSClientConfig.InsecureSkipVerify = true
		}

		if cfg.PiHole.CAFile != "" {
			cadata, err := ioutil.ReadFile(cfg.PiHole.CAFile)
			if err != nil {
				return result, err
			}

			cacerts := x509.NewCertPool()
			if !cacerts.AppendCertsFromPEM(cadata) {
				return result, fmt.Errorf("Can't append CA data to CA pool")
			}

			transp.TLSClientConfig.RootCAs = cacerts
		}
	}

	cl := &http.Client{
		Transport: transp,
		Timeout:   cfg.PiHole.timeout,
	}

	if !cfg.PiHole.FollowRedirect {
		cl.CheckRedirect = func(http_request *http.Request, http_via []*http.Request) error { return http.ErrUseLastResponse }
	}

	piurl = cfg.PiHole.URL + "?" + stat

	if cfg.PiHole.AuthHash != "" {
		piurl += "&auth=" + cfg.PiHole.AuthHash
	}

	request, err := http.NewRequest("GET", piurl, nil)
	if err != nil {
		return result, err
	}

	// always consume HTTP request body
	defer func() {
		if request.Body != nil {
			ioutil.ReadAll(request.Body)
			request.Body.Close()
		}
	}()

	request.Header.Set("User-Agent", userAgent)

	/*
	   "A man is not dead while his name is still spoken."
	   - Going Postal, Chapter 4 prologue
	*/
	request.Header.Set("X-Clacks-Overhead", "GNU Terry Pratchett")

	// close TCP session
	request.Close = true

	response, err := cl.Do(request)
	if err != nil {
		return result, err
	}

	// always consume reply
	defer func() {
		ioutil.ReadAll(response.Body)
		response.Body.Close()
	}()

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.Header = response.Header
	result.Content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	return result, nil
}
