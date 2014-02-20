/*
strip_frame_sources takes a DataBurst on stdin and writes a
DataSourceBurst to stdout.
*/
package main

import (
	"bytes"
	"github.com/anchor/bletchley/dataframe"
	"github.com/anchor/chevalier"
	"log"
	"os"
)

func main() {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	burst, err := dataframe.UnmarshalDataBurst(buf.Bytes())
	if err != nil {
		log.Fatal("Could not unmarshal DataBurst: %v", err)
	}
	sources := make([]*chevalier.DataSource, len(burst.Frames))
	for i, frame := range burst.Frames {
		sources[i] = chevalier.DataFrameSource(frame)
	}
	sourceBurst := chevalier.BuildSourceBurst(sources)
	if err != nil {
		log.Fatal("Could not unmarshal source: %v", err)
	}
	out, err := chevalier.MarshalSourceBurst(sourceBurst)
	if err != nil {
		log.Fatal("Could not marshal SourceBurst: %v", err)
	}
	os.Stdout.Write(out)
}
