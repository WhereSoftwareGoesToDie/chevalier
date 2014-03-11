package main

import (
	"flag"
	"github.com/anchor/chevalier"
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func handleErrors(w *chevalier.ElasticsearchWriter) {
	ch := w.GetErrorChan()
	for errBuf := range ch {
		log.Println(errBuf.Err)
	}
}

func main() {
	esHost := flag.String("host", "localhost", "Elasticsearch host to connect to")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s <origin> [args]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	origin := flag.Arg(0)
	writer := chevalier.NewElasticsearchWriter(*esHost, 1, 60, "chevalier", "datasource")
	reader := io.Reader(os.Stdin)
	packet, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("Could not read from stdin: %v", err)
	}
	burst, err := chevalier.UnmarshalSourceBurst(packet)
	if err != nil {
		log.Fatal("Could not unmarshal source: %v", err)
	}
	for _, source := range burst.Sources {
		err = writer.Write(origin, source)
		if err != nil {
			log.Println("Writer error: %v", err)
		}
		go handleErrors(writer)
	}
	writer.Shutdown()
}
