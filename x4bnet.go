package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func x4bNetList(wtr *mmdbwriter.Tree) {
	urls := map[string]string{
		"https://raw.githubusercontent.com/X4BNet/lists_vpn/main/output/vpn/ipv4.txt":        "is_vpn",           // 43.26
		"https://raw.githubusercontent.com/X4BNet/lists_vpn/main/output/datacenter/ipv4.txt": "is_datacenter",    //361.61
		"https://raw.githubusercontent.com/X4BNet/lists_torexit/main/ipv4.txt":               "is_tor_exit",      //26.85
		"https://raw.githubusercontent.com/X4BNet/lists_stopforumspam/main/ipv4.txt":         "is_stopForumSpam", //69.79
	}

	botUrls := map[string]string{
		"https://raw.githubusercontent.com/X4BNet/lists_searchengine/main/outputs/ahrefs.txt": "ahref",     //1.24
		"https://raw.githubusercontent.com/X4BNet/lists_searchengine/main/outputs/bing.txt":   "bing",      //0.99
		"https://raw.githubusercontent.com/X4BNet/lists_searchengine/main/outputs/google.txt": "google",    //0.62
		"https://raw.githubusercontent.com/X4BNet/lists_searchengine/main/outputs/mj12.txt":   "mj12",      //51.71
		"https://raw.githubusercontent.com/X4BNet/lists_searchengine/main/outputs/yahoo.txt":  "yahoo",     //3.03
		"https://raw.githubusercontent.com/X4BNet/lists_bots/main/outputs/operamini.txt":      "operamini", //2.91
		"https://raw.githubusercontent.com/X4BNet/lists_bots/main/outputs/semrush.txt":        "semrush",   //58.38
	}
	for url, key := range urls {
		if err := extractor(url, key, wtr, boolGenerator); err != nil {
			log.Fatal(err)
		}
	}
	for url, key := range botUrls {
		if err := extractor(url, key, wtr, botInfoGenerator); err != nil {
			log.Fatal(err)
		}
	}

	fh, err := os.Create("resources/IP-external-dir.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = wtr.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

func extractor(url, key string, wtr *mmdbwriter.Tree, generator func(string) mmdbtype.DataType) error {
	// Fetch the content from the URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetching the URL:: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch the content. Status code: %d", resp.StatusCode)
	}

	// Parse the response content line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// Each line contains an IPv4 address
		ipv4Address := scanner.Text()

		_, ipNet, err := net.ParseCIDR(ipv4Address)
		if err != nil {
			log.Fatal(err)
		}
		if err := wtr.InsertFunc(ipNet, inserter.TopLevelMergeWith(generator(key))); err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	return nil
}

func boolGenerator(key string) mmdbtype.DataType {
	return mmdbtype.Map{mmdbtype.String(key): mmdbtype.Bool(true)}
}
func botInfoGenerator(key string) mmdbtype.DataType {
	return mmdbtype.Map{"botInfo": mmdbtype.Map{
		"is_bot": mmdbtype.Bool(true),
		"name":   mmdbtype.String(key),
	}}
}
