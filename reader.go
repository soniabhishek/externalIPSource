package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

func reader() {
	db, err := maxminddb.Open("resources/IP-external-dir.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP("34.73.206.133")

	var record any
	err = db.Lookup(ip, &record)
	if err != nil {
		log.Panic(err)
	}
	c, _ := json.Marshal(record)
	fmt.Printf("%v,\n %s", record, string(c))
}
