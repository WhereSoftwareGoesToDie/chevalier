package main

import (
	"flag"
	"fmt"
	"github.com/anchor/chevalier"
	zmq "github.com/pebbe/zmq4"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func queryES(origin string, req *chevalier.SourceRequest, host string) {
	engine := chevalier.NewQueryEngine(host, "chevalier", "datasource")
	results, err := engine.GetSources(origin, req)
	if err != nil {
		log.Println("Search error: %v", err)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, source := range results.GetSources() {
		fmt.Println(source)
	}
}

func queryChevalier(origin string, req *chevalier.SourceRequest, endpoint string) {
	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatal(err)
	}
	err = sock.Connect(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	packet, err := chevalier.MarshalSourceRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	_, err = sock.SendMessage(origin, packet, 0)
	response, err := sock.RecvBytes(0)
	if err != nil {
		log.Fatal(err)
	}
	burst, err := chevalier.UnmarshalSourceBurst(response)
	if err != nil {
		log.Fatal(err)
	}
	for _, source := range burst.GetSources() {
		fmt.Println(source)
	}
}

func main() {
	esHost := flag.String("host", "localhost", "Elasticsearch host to connect to")
	protobuf := flag.Bool("protobuf", false, "Read a SourceRequest from stdin rather than accepting field:value pairs on the command line.")
	es := flag.Bool("es", false, "Read from Elasticsearch directly rather than chevalier.")
	startPage := flag.Int("start-page", 0, "Obtain results from this page.")
	pageSize := flag.Int("page-size", 0, "Number of results per page.")
	endpoint := flag.String("endpoint", "tcp://127.0.0.1:6283", "Chevalier endpoint (as a ZMQ URI).")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s <origin> <field:value> [field:value ...] [args]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	origin := flag.Arg(0)
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
		for i, arg := range flag.Args()[1:] {
			pair := strings.Split(arg, ":")
			if len(pair) < 2 {
				log.Fatal("Could not parse %v: must be a 'field:value' pair.")
			}
			tags[i] = chevalier.NewSourceRequestTag(pair[0], pair[1])
		}
		req = chevalier.NewSourceRequest(tags)
		if *startPage > 0 {
			page := int64(*startPage)
			req.StartPage = &page
		}
		if *pageSize > 0 {
			size := int64(*pageSize)
			req.SourcesPerPage = &size
		}
	}
	if *es {
		queryES(origin, req, *esHost)
	} else {
		queryChevalier(origin, req, *endpoint)
	}
}
