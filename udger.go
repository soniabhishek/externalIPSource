package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

const UdgerSecret = "XXXXXXXXXXXXX"

func udgerList(tree *mmdbwriter.Tree) error {
	urls := map[string]string{
		"https://data.udger.com/" + UdgerSecret + "/anonymizing_vpn_service.csv":  "is_anonymizing_vpn",
		"https://data.udger.com/" + UdgerSecret + "/fake_crawler.csv":             "is_fake_crawler",
		"https://data.udger.com/" + UdgerSecret + "/tor_exit_node.csv":            "tor_exit_node",
		"https://data.udger.com/" + UdgerSecret + "/web_scraper.csv":              "web_scraper",
		"https://data.udger.com/" + UdgerSecret + "/known_attack_source.csv":      "known_attack_source",
		"https://data.udger.com/" + UdgerSecret + "/known_attack_source_mail.csv": "known_attack_source_mail",
		"https://data.udger.com/" + UdgerSecret + "/known_attack_source_ssh.csv":  "known_attack_source_ssh",
		"https://data.udger.com/" + UdgerSecret + "/public_cgi_proxy.csv":         "public_cgi_proxy",
		"https://data.udger.com/" + UdgerSecret + "/public_web_proxy.csv":         "public_web_proxy",
	}

	longUrls := map[string]string{
		"https://data.udger.com/" + UdgerSecret + "/datacenter.csv": "is_datacenter",
	}

	hugeUrls := map[string]string{
		"https://data.udger.com/" + UdgerSecret + "/botIP.csv": "is_bot",
	}

	for url, key := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("unable to get on url : %s, err: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", resp.Status)
		}
		if err := transform(tree, resp.Body, key); err != nil {
			return fmt.Errorf("transform csv: %s, err: %w", url, err)
		}
	}

	for url, key := range longUrls {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("unable to get on url : %s, err: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", resp.Status)
		}
		if err := parseDataCentersFile(tree, resp.Body, key); err != nil {
			return fmt.Errorf("transform csv: %s, err: %w", url, err)
		}
	}

	for url, key := range hugeUrls {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("unable to get on url : %s, err: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", resp.Status)
		}
		if err := transformBotFile(tree, resp.Body, key); err != nil {
			return fmt.Errorf("transform csv: %s, err: %w", url, err)
		}
	}
	return nil
}

func transform(tree *mmdbwriter.Tree, raw io.Reader, key string) error {
	reader := csv.NewReader(raw)
	for {
		rec, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("transform udger csv, %w", err)
		}

		ip := net.ParseIP(rec[0])
		if ip == nil {
			log.Printf("invalid ip: %s", rec[0])
		}

		if err := tree.InsertRangeFunc(
			ip,
			ip,
			inserter.TopLevelMergeWith(mmdbtype.Map{mmdbtype.String(key): mmdbtype.Bool(true)}),
		); err != nil {
			log.Printf("insert to tree: %s", err.Error())
			continue
		}
	}
	return nil
}

func parseDataCentersFile(wtr *mmdbwriter.Tree, raw io.Reader, key string) error {
	reader := csv.NewReader(raw)

	for {
		rec, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read data center csv, %w", err)
		}
		ipStart := net.ParseIP(rec[2])
		ipEnd := net.ParseIP(rec[3])
		if ipStart == nil || ipEnd == nil {
			log.Printf("invalid range: %s, %s", rec[2], rec[3])
		}
		if err := wtr.InsertRangeFunc(
			ipStart,
			ipEnd,
			inserter.TopLevelMergeWith(mmdbtype.Map{
				mmdbtype.String(key): mmdbtype.Bool(true),
				"datacenter_name":    mmdbtype.String(rec[0]),
			})); err != nil {
			log.Printf("insert to tree: %s", err.Error())
			continue
		}
	}

	return nil
}

func transformBotFile(wtr *mmdbwriter.Tree, raw io.Reader, key string) error {
	reader := csv.NewReader(raw)

	for {
		rec, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read data center csv, %w", err)
		}
		ip := net.ParseIP(rec[1])
		if ip == nil {
			log.Printf("invalid ip: %s", rec[1])
		}
		if err := wtr.InsertRangeFunc(
			ip,
			ip,
			inserter.TopLevelMergeWith(mmdbtype.Map{
				mmdbtype.String(key): mmdbtype.Bool(true),
				"bot_name":           mmdbtype.String(rec[5]),
				"bot_category":       mmdbtype.String("string"),
			})); err != nil {
			log.Printf("insert to tree: %s", err.Error())
			continue
		}
	}

	return nil
}
