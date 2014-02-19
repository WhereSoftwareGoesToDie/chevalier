/* 
strip_frame_sources takes a DataBurst on stdin and writes a
DataSourceBurst to stdout.
*/
package main

import (
	"github.com/anchor/bletchley/dataframe"
	"github.com/anchor/chevalier"
	"os"
	"log"
	"bytes"
)

func frameSource(frame *dataframe.DataFrame) *chevalier.DataSource {
	source := new(chevalier.DataSource)
	source.Source = make([]*chevalier.DataSource_Tag, len(frame.Source))
	for i, tag := range frame.Source {
		source.Source[i] = new(chevalier.DataSource_Tag)
		source.Source[i].Field = tag.Field
		source.Source[i].Value = tag.Value
	}
	return source
}

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
		sources[i] = frameSource(frame)
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
