package main

import (
	"flag"
	"time"

	"github.com/anchor/chevalier"
	"github.com/anchor/zmqutil"
	"github.com/fractalcat/nagiosplugin"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	check := nagiosplugin.NewCheck()
	defer check.Finish()
	endpoint := flag.String("endpoint", "tcp://127.0.0.1:6284", "Chevalier status endpoint (as a ZMQ URI).")
	flag.Parse()
	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		check.Criticalf("Could not create socket: %v", err)
	}
	err = sock.Connect(*endpoint)
	if err != nil {
		check.Criticalf("Could connect to Chevalier endpoint %v: %v", *endpoint, err)
	}
	_, err = zmqutil.RetrySend(sock, "", 0)
	if err != nil {
		check.Criticalf("Could not request status from Chevalier: %v", err)
	}
	p, err := zmqutil.RetryRecvBytes(sock, 0)
	if err != nil {
		check.Criticalf("Could not receive status from Chevalier: %v", err)
	}
	status, err := chevalier.UnmarshalStatusResponse(p)
	if err != nil {
		check.Criticalf("Could not unmarshal Chevalier status message: %v", err)
	}
	if status.GetErrors() != nil {
		es := status.GetErrors()
		for _, e := range es {
			check.AddResult(nagiosplugin.WARNING, e)
		}
	}
	now := time.Now().UnixNano()
	originCount := 0
	totalSources := uint64(0)
	for _, o := range status.Origins {
		originCount++
		totalSources += o.GetSources()
		lastUpdated := float64(uint64(now) - o.GetLastUpdated())
		lastUpdated /= 1000000000.0 // nanoseconds -> seconds
		check.AddPerfDatum(*o.Origin + "_sources", "", float64(o.GetSources()))
		check.AddPerfDatum(*o.Origin + "_age", "s", lastUpdated)
	}
	check.AddPerfDatum("total_sources", "", float64(totalSources))
	check.AddResult(nagiosplugin.OK, "Chevalier seems alive")
}
