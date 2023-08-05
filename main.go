package main

import (
	"log"

	"github.com/maxmind/mmdbwriter"
)

func main() {
	//fromScratch()
	//aggregate()
	reader()
}

func aggregate() {
	wtr, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "IP-External-Source-Directory",
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := udgerList(wtr); err != nil {
		log.Fatal(err)
	}
	x4bNetList(wtr)
}
