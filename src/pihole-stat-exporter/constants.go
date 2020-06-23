package main

const name = "pihole-stat-exporter"
const version = "1.0.0-20200623"

var userAgent = name + "/" + version

const defaultExporterURL = "http://127.0.0.1:64711"
const defaultPrometheusPath = "/metrics"
const defaultInfluxDataPath = "/influx"

const versionText = `%s version %s
Copyright (C) 2020 by Andreas Maus <maus@ypbind.de>
This program comes with ABSOLUTELY NO WARRANTY.

%s is distributed under the Terms of the GNU General
Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)

Build with go version: %s

`

const helpText = `Usage: %s --config=<cfg> [--help] [--version]
    --config=<cfg>  Path to the configuration file
                    This parameter is mandatory

    --help          This help text

    --version       Show version information

`
