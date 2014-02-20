package main

import (
	"flag"
	"github.com/anchor/chevalier"
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	esHost := flag.String("host", "localhost", "Elasticsearch host to connect to")
	protobuf := flag.Bool("protobuf", false, "Read a SourceRequest from stdin rather than accepting field:value pairs on the command line.")
	flag.Parse()
	engine := chevalier.NewQueryEngine(*esHost, "chevalier", "datasource")
	var req *chevalier.SourceRequest
	if *protobuf {
		reader := io.Reader(os.Stdin)
		packet, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("Could not read from stdin: %v", err)
		}
		req, err = chevalier.UnmarshalSourceRequest(packet)
		if err != nil {
			log.Fatal("Could not unmarshal request: %v", err)
		}
	} else {
		tags := make([]*chevalier.SourceRequest_Tag, flag.NArg())
		for i, arg := range flag.Args() {
			pair := strings.Split(arg, ":")
			if len(pair) < 2 {
				log.Fatal("Could not parse %v: must be a 'field:value' pair.")
			}
			tags[i] = chevalier.NewSourceRequestTag(pair[0], pair[1])
		}
		req = chevalier.NewSourceRequest(tags)
	}
	results, err := engine.RunSourceRequest(req)
	if err != nil {
		log.Println("Search error: %v", err)
	}
	sources := chevalier.FmtResult(results)
	if err != nil {
		log.Fatal(err)
	}
	for _, source := range sources {
		fmt.Println(source)
	}
}
