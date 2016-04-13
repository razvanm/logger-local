package main

import (
	"flag"
	"log"

	"fmt"
	"v.io/v23"
	"v.io/v23/syncbase"
	_ "v.io/x/ref/runtime/factories/generic"
)

var (
	sbservice = flag.String("syncbase", "/localhost:8151/syncbase/server", "Syncbase service to use")
)

func main() {
	ctx, shutdown := v23.Init()
	defer shutdown()

	service := syncbase.NewService(*sbservice)
	db := service.Database(ctx, "default", nil)
	collection := db.Collection(ctx, "logger")

	resumeMarker, err := db.GetResumeMarker(ctx)
	if err != nil {
		log.Fatalf("failed getting the resume maker: %v", err)
	}
	watchStream, err := db.Watch(ctx, collection.Id(), "", resumeMarker)
	if err != nil {
		log.Fatalf("failed start watching: %v", err)
	}

	for watchStream.Advance() {
		change := watchStream.Change()
		switch change.ChangeType {
		case syncbase.PutChange:
			var text string
			change.Value(&text)
			fmt.Printf("Changed key: %q Value: %q\n", change.Row, text)
		case syncbase.DeleteChange:
			fmt.Printf("Deleted key: %q\n", change.Row)
		}

	}
	fmt.Printf("WatchStream error: %v", watchStream.Err())
}
