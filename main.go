package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mtulio/azion-exporter/src/azion"
	"github.com/prometheus/common/log"
)

var (
	fEmail    *string
	fPassword *string
)

// usage returns the command line usage sample.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	fEmail = flag.String("azion.email", "", "API email address to get Authorization token")
	fPassword = flag.String("azion.password", "", "API password to get Authorization token")
	flag.Usage = usage
	flag.Parse()

	if *fEmail == "" {
		*fEmail = os.Getenv("AZION_EMAIL")
	}

	if *fPassword == "" {
		*fPassword = os.Getenv("AZION_PASSWORD")
	}

}

func main() {
	log.Infoln("Starting exporter ")
	c := azion.NewClient(*fEmail, *fPassword)
	meta, err := c.Analytics.GetMatadata()
	if err != nil {
		panic(err)
	}
	fmt.Println(meta)

	// var metric *azion.Analytics.MetricResponse
	metric, err2 := c.Analytics.GetCDRequestsTotal()
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(metric)
}
