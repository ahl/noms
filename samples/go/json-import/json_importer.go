// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/util/jsontonoms"
	"github.com/attic-labs/noms/go/util/progressreader"
	"github.com/attic-labs/noms/go/util/status"
	"github.com/dustin/go-humanize"
	flag "github.com/juju/gnuflag"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <url> <dataset>\n", os.Args[0])
		flag.PrintDefaults()
	}

	spec.RegisterDatabaseFlags(flag.CommandLine)
	flag.Parse(true)

	if len(flag.Args()) != 2 {
		d.CheckError(errors.New("expected url and dataset flags"))
	}

	ds, err := spec.GetDataset(flag.Arg(1))
	d.CheckError(err)

	url := flag.Arg(0)
	if url == "" {
		flag.Usage()
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching %s: %+v\n", url, err)
	} else if res.StatusCode != 200 {
		log.Fatalf("Error fetching %s: %s\n", url, res.Status)
	}
	defer res.Body.Close()

	var jsonObject interface{}
	start := time.Now()
	r := progressreader.New(res.Body, func(seen uint64) {
		elapsed := time.Since(start).Seconds()
		rate := uint64(float64(seen) / elapsed)
		status.Printf("%s decoded in %ds (%s/s)...", humanize.Bytes(seen), int(elapsed), humanize.Bytes(rate))
	})
	err = json.NewDecoder(r).Decode(&jsonObject)
	if err != nil {
		log.Fatalln("Error decoding JSON: ", err)
	}
	status.Done()

	_, err = ds.CommitValue(jsontonoms.NomsValueFromDecodedJSON(jsonObject, true))
	d.PanicIfError(err)
	ds.Database().Close()
}
