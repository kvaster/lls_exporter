package main

import (
	"flag"
	"github.com/apex/log"
	"github.com/kvaster/apexutils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"lls_exporter/lls"
	"os"
	"os/signal"
	"syscall"
)

var configFile = flag.String("config", "/etc/lls-exporter/lls.yml", "config file")

func main() {
	flag.Parse()
	apexutils.ParseFlags()

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.WithField("error", err).Fatal("error reading config file")
	}

	var config lls.SensorsConfig
	err = yaml.UnmarshalStrict(data, &config)
	if err != nil {
		log.WithField("error", err).Fatal("error parsing config file")
	}

	log.Info("starting lls exporter")

	exporter := lls.New()

	go func() {
		stopChan := make(chan os.Signal)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		<-stopChan

		exporter.Shutdown()
	}()

	err = exporter.Serve(config)
	if err != nil {
		log.WithField("error", err).Fatal("serve error")
	}

	log.Info("lls exported stopped")
}
