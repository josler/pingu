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

var filepathVar string
var limitVar int

func main() {
	flag.Parse()
	file, _ := os.Open("resources/pingu-config.json")
	c, err := core.ParseConfig(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var reader io.Reader
	if filepathVar != "" {
		fmt.Printf("Reading: %s\n\n", filepathVar)
		file, err := os.Open(filepathVar)
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
}

func init() {
	flag.StringVar(&filepathVar, "f", "", "filepath to parse")
	flag.IntVar(&limitVar, "l", -1, "number of items to limit report to")
}
