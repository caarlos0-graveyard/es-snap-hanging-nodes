package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

var esURL = flag.String("es.url", "http://localhost:9200", "ElasticSearch URL")

func main() {
	flag.Parse()

	resp, err := http.Get(*esURL + "/_snapshot/_status?pretty")
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("got a non-ok status:", resp.StatusCode)
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var r Result
	if err := json.Unmarshal(bts, &r); err != nil {
		log.Fatalln(err)
	}

	indices := map[string][]string{}
	for _, snap := range r.Snapshots {
		for indexName, index := range snap.Indices {
			for sharNumber, stat := range index.Shards {
				if stat.Stage == "INIT" {
					indices[indexName] = append(indices[indexName], sharNumber)
					break
				}
			}
		}
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()

	for idx, shards := range indices {
		fmt.Fprintf(w, "%s\t%s\n", idx, strings.Join(shards, ","))
	}
}

type Result struct {
	Snapshots []struct {
		Indices map[string]struct {
			Shards map[string]struct {
				Stage string `json:"stage"`
			} `json:"shards"`
		} `json:"indices"`
	} `json:"snapshots"`
}
