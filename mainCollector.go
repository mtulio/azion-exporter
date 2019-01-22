package main

import (
	"github.com/mtulio/azion-exporter/src/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func initPromCollector() error {
	var err error
	err = nil
	if cfg.prom == nil {
		cfg.prom = new(globalProm)
	}

	cfg.prom.Collector, err = collector.NewCollectorMaster(cfg.azionClient, cfg.metricsName...)
	if err != nil {
		log.Warnln("Init Prom: Couldn't create collector: ", err)
		return err
	}

	cfg.prom.Registry = prometheus.NewRegistry()
	err = cfg.prom.Registry.Register(cfg.prom.Collector)
	if err != nil {
		log.Errorln("Init Prom: Couldn't register collector:", err)
		return err
	}

	cfg.prom.Gatherers = &prometheus.Gatherers{
		prometheus.DefaultGatherer,
		cfg.prom.Registry,
	}
	return nil
}
