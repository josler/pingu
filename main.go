package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/josler/pingu/core"
	"github.com/josler/pingu/processor"
	"github.com/josler/pingu/tailer"
)

var logFilepathVar string
var configFilepathVar string
var limitVar int

func main() {
	flag.Parse()
	file, err := os.Open(configFilepathVar)
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := core.ParseConfig(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var reader io.Reader
	if logFilepathVar != "" {
		fmt.Printf("Reading: %s\n\n", logFilepathVar)
		file, err := os.Open(logFilepathVar)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		reader = file
	} else {
		reader = os.Stdin
	}

	p := tailer.Parser{}
	p.Parse(reader)

	rule := processor.ParseRules(c.Rules())
	f := processor.ParseCalculation(c.Calculation())
	filter := processor.NewFilter(rule, p.Out())
	grouping := processor.NewGrouping(
		c.GroupBy(),
		f,
		filter.Out(),
	)

	grouping.Start()
	fmt.Println(grouping.Report(limitVar))
	fmt.Println(grouping.ReportBuckets(limitVar))
}

func init() {
	flag.StringVar(&logFilepathVar, "f", "", "filepath to logs to parse")
	flag.StringVar(&configFilepathVar, "c", "resources/pingu-config.json", "filepath to config")
	flag.IntVar(&limitVar, "l", -1, "number of items to limit report to")
}
