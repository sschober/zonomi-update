package main

import (
	"flag"
	"fmt"
	"github.com/vharitonsky/iniflags"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
)

const APP_VERSION = "0.1"
const INVALID string = "invalid"

const zonomiUpdateURL string = "https://zonomi.com/app/dns/dyndns.jsp?host=%s&api_key=%s"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var debugFlag *bool = flag.Bool("debug", false, "Enable dbugging output")
var disableSyslogFlag *bool = flag.Bool("disable-syslog", false, "Disable syslog logging")

var domainFlag *string = flag.String("domain", INVALID, "Domain to update")
var apiKeyFlag *string = flag.String("apikey", INVALID, "Zonomi API key")

var syslogger *syslog.Writer
var syslogEnabled bool

func syslogFatal(msg string) {
	if syslogEnabled {
		if err := syslogger.Err(msg); err != nil {
			log.Fatalf("logging to syslog failed: %v", err)
		}
	}
}

func syslogInfo(msg string) {
	if syslogEnabled {
		if err := syslogger.Info(msg); err != nil {
			log.Fatalf("logging to syslog failed: %v", err)
		}
	}
}

func validateFlags() {
	if INVALID == *apiKeyFlag {
		msg := "api key not set. exiting."
		syslogFatal(msg)
		log.Fatalf(msg)
	}
	if INVALID == *domainFlag {
		msg := "domain not set. exiting."
		syslogFatal(msg)
		log.Fatalf(msg)
	}
}
func main() {
	iniflags.SetConfigFile(os.Getenv("HOME") + "/.zonomi-update")
	iniflags.Parse()

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		return
	}

	if !*debugFlag {
		log.SetOutput(ioutil.Discard)
	}

	validateFlags()

	syslogEnabled = !*disableSyslogFlag

	if syslogEnabled {
		sl, err := syslog.New(syslog.LOG_NOTICE, "zonomi-update")
		if err != nil {
			log.Fatalf("Could not initialize syslog logging: %s", err.Error())
		}
		syslogger = sl
		defer syslogger.Close()
	}

	url := fmt.Sprintf(zonomiUpdateURL, *domainFlag, *apiKeyFlag)
	log.Printf("update url: %s", url)

	resp, err := http.Get(url)

	if err != nil {
		syslogFatal(fmt.Sprintf("updating domain '%s' failed: %v", *domainFlag, err))
		log.Fatalf("updating domain '%s' failed: %v", *domainFlag, err)
	}

	if resp.StatusCode != 200 {
		syslogFatal(fmt.Sprintf("updating domain '%s' failed: status: %s", *domainFlag, resp.Status))
		log.Fatalf("updating domain '%s' failed: response: %v", *domainFlag, resp)
	}

	log.Printf("domain '%s' updated:\n%v", *domainFlag, resp)
	syslogInfo(fmt.Sprintf("domain '%s' updated: %s", *domainFlag, resp.Status))
}
