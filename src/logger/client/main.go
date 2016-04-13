package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"v.io/v23"
	wire "v.io/v23/services/syncbase"
	"v.io/v23/syncbase"
	_ "v.io/x/ref/runtime/factories/generic"
)

var (
	sbservice = flag.String("syncbase", "/localhost:8151/syncbase/client", "Syncbase service to use")
	duration  = flag.Duration("duration", 2*time.Second, "How long to run")
	sleep     = flag.Duration("sleep", 250*time.Millisecond, "How much to sleep between two messages")
)

func main() {
	ctx, shutdown := v23.Init()
	defer shutdown()

	service := syncbase.NewService(*sbservice)
	db := service.Database(ctx, "default", nil)
	collection := db.Collection(ctx, "logger")

	tEnd := time.NewTimer(*duration)
	for {
		if err := syncbase.RunInBatch(ctx, db, wire.BatchOptions{}, func(db syncbase.BatchDatabase) error {
			key := fmt.Sprintf("%d", time.Now().UnixNano())
			value := fmt.Sprintf("text-%010v", rand.Int31())
			fmt.Printf("Inserting (%s, %s) ...\n", key, value)
			if err := collection.Put(ctx, key, value); err != nil {
				return err
			}
			return nil
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		tSleep := time.NewTimer(*sleep)
		select {
		case <-tEnd.C:
			return
		case <-tSleep.C:
			break
		}
	}
}
