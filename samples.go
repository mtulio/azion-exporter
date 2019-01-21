package main

import (
	"fmt"

	"github.com/mtulio/azion-exporter/src/azion"
)

func sampleGetMetadata(c *azion.Client) {
	meta, err := c.Analytics.GetMatadata()
	if err != nil {
		panic(err)
	}
	fmt.Println(meta)
	return
}

func sampleGetMetricProdCDDimension(c *azion.Client) {
	metric, err := c.Analytics.GetMetricDimensionProdCD("requests", "total", "date_from=last-hour")
	if err != nil {
		panic(err)
	}
	fmt.Println(metric)
}
