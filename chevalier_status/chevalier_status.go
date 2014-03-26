package main

import (
	"fmt"
	"flag"
	"log"
	"os"

	"github.com/anchor/chevalier"
	"github.com/anchor/zmqutil"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	endpoint := flag.String("endpoint", "tcp://127.0.0.1:6284", "Chevalier status endpoint (as a ZMQ URI).")
	flag.Parse()
	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatal(err)
	}
	err = sock.Connect(*endpoint)
	if err != nil {
		log.Fatal(err)
	}
	_, err = zmqutil.RetrySend(sock, "", 0)
	if err != nil {
		log.Fatal(err)
	}
	p, err := zmqutil.RetryRecvBytes(sock, 0)
	if err != nil {
		log.Fatal(err)
	}
	status, err := chevalier.UnmarshalStatusResponse(p)
	if err != nil {
		log.Fatal(err)
	}
	out, err := status.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(out)
	fmt.Println()
}
