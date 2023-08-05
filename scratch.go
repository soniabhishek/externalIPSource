package main

import (
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func fromScratch() {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "IP-External-Source-Directory",
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	_, ipNet, err := net.ParseCIDR("91.211.116.0")
	if err != nil {
		log.Fatal(err)
	}
	if err := writer.Insert(ipNet, mmdbtype.Map{"dataCenter": mmdbtype.String("0X2A_DATACENTER")}); err != nil {
		log.Fatal(err)
	}
	fh, err := os.Create("resources/IP-external-dir.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
