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

func queryES(req *chevalier.SourceRequest, host string) {
	engine := chevalier.NewQueryEngine(host, "chevalier", "datasource")
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

func queryChevalier(req *chevalier.SourceRequest, endpoint string) {
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
	_, err = sock.SendBytes(packet, 0)
	response, err := sock.RecvMessageBytes(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}

func main() {
	esHost := flag.String("host", "localhost", "Elasticsearch host to connect to")
	protobuf := flag.Bool("protobuf", false, "Read a SourceRequest from stdin rather than accepting field:value pairs on the command line.")
	es := flag.Bool("es", false, "Read from Elasticsearch directly rather than chevalier.")
	endpoint := flag.String("endpoint", "tcp://127.0.0.1:6283", "Chevalier endpoint (as a ZMQ URI).")
	flag.Parse()
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
	if *es {
		queryES(req, *esHost)
	} else {
		queryChevalier(req, *endpoint)
	}
}
