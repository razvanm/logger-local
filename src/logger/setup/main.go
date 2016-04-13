package main

import (
	"flag"
	"fmt"
	"log"

	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/naming"
	"v.io/v23/security"
	"v.io/v23/security/access"
	wire "v.io/v23/services/syncbase"
	"v.io/v23/syncbase"
	_ "v.io/x/ref/runtime/factories/generic"
)

var (
	sgMountPoint = flag.String("mountpoint", "", "Syncgroup mount point")
)

func setup(ctx *context.T, name string) {
	service := syncbase.NewService(name)

	db := service.Database(ctx, "default", nil)
	if exists, err := db.Exists(ctx); err != nil {
		log.Fatalf("failed checking for db %q: %v", db.FullName(), err)
	} else if !exists {
		log.Printf("Creating database %q...", db.FullName())
		if err := db.Create(ctx, nil); err != nil {
			log.Fatalf("failed creating database %q: %v", db.FullName(), err)
		}
	}

	collection := db.Collection(ctx, "logger")
	if exists, err := collection.Exists(ctx); err != nil {
		log.Fatalf("failed checking for collection %q: %v", collection.FullName(), err)
	} else if !exists {
		log.Printf("Creating collection %q...", collection.FullName())
		if err := collection.Create(ctx, nil); err != nil {
			log.Fatalf("failed creating collection %q: %v", collection, err)
		}
	}
}

func createSyncgroup(ctx *context.T, name, sgName string) {
	service := syncbase.NewService(name)
	db := service.Database(ctx, "default", nil)
	collection := db.Collection(ctx, "logger")
	sg := db.Syncgroup(sgName)

	allAccess := access.AccessList{In: []security.BlessingPattern{"..."}}
	spec := wire.SyncgroupSpec{
		Perms: access.Permissions{
			"Admin":   allAccess,
			"Write":   allAccess,
			"Read":    allAccess,
			"Resolve": allAccess,
			"Debug":   allAccess,
		},
		Prefixes:    []wire.CollectionRow{wire.CollectionRow{CollectionId: collection.Id()}},
		MountTables: []string{*sgMountPoint},
	}
	err := sg.Create(ctx, spec, wire.SyncgroupMemberInfo{})
	if err != nil {
		fmt.Printf("Failed to create the syncgroup: %v\n", err)
	} else {
		fmt.Printf("Created syncgroup %q: %+v\n", sgName, spec)
	}
}

func joinSyncgroup(ctx *context.T, name, sgName string) {
	service := syncbase.NewService(name)
	db := service.Database(ctx, "default", nil)
	sg := db.Syncgroup(sgName)
	spec, err := sg.Join(ctx, wire.SyncgroupMemberInfo{})
	if err != nil {
		fmt.Printf("Failed to join the syncgroup: %v\n", err)
	} else {
		fmt.Printf("Joined syncgroup %q: %+v\n", sgName, spec)
	}
}

func main() {
	ctx, shutdown := v23.Init()
	defer shutdown()

	sgName := naming.Join(flag.Args()[0], "%%sync", "logger")
	for _, name := range flag.Args()[:1] {
		fmt.Printf("Setting up %q\n", name)
		setup(ctx, name)
		createSyncgroup(ctx, name, sgName)
	}

	for _, name := range flag.Args()[1:] {
		fmt.Printf("Setting up %q\n", name)
		setup(ctx, name)
		joinSyncgroup(ctx, name, sgName)
	}
}
