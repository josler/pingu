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
		file = getFileAtPos(0, logFilepathVar)
		reader = file
		defer file.Close()
	} else {
		reader = os.Stdin
	}
	p := tailer.Parser{}
	p.Parse(reader)

	query := processor.NewQuery(&p, limitVar, c)
	query.Start()
}

func getFileAtPos(position int64, filename string) *os.File {
	fmt.Printf("Reading: %s\n\n", filename)
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	file.Seek(position, 0)
	return file
}

func init() {
	flag.StringVar(&logFilepathVar, "f", "", "filepath to logs to parse")
	flag.StringVar(&configFilepathVar, "c", "resources/pingu-config.json", "filepath to config")
	flag.IntVar(&limitVar, "l", -1, "number of items to limit report to")
}
