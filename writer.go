package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func writer() {
	wtr, err := rawWriter("resources/GeoIP2-City.mmdb", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wtr)

	_, sreNet, err := net.ParseCIDR("56.0.2.2/32")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sreNet)
	sreData := mmdbtype.Map{
		"Fingerprint.botInfo":    mmdbtype.String("Google"),
		"Fingerprint.datacenter": mmdbtype.String("Google"),
	}
	if err := wtr.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
		log.Fatal(err)
	}
	_, sreNet, err = net.ParseCIDR("56.0.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sreNet)
	sreData = mmdbtype.Map{
		"Fingerprint.XXXXX": mmdbtype.String("XXXXX"),
	}
	if err := wtr.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
		log.Fatal(err)
	}
	fh, err := os.Create("resources/GeoData-Updated.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = wtr.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

func rawWriter(path, name string) (*mmdbwriter.Tree, error) {
	wtr, err := mmdbwriter.Load(path, mmdbwriter.Options{
		DatabaseType: name,
		RecordSize:   24,
	})
	return wtr, err
}
