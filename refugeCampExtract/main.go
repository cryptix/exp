package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// from: https://www.hamburg.de/fluechtlinge-unterbringung-standorte/
var urls = []string{
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4372596/unterbringung-altona/",
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373118/unterbringung-bergedorf/",
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373120/unterbringung-eimsbuettel/",
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373122/unterbringung-harburg/",
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373128/unterbringung-mitte/",
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373126/unterbringung-nord/",
	"http://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373124/unterbringung-wandsbek/",
}

func main() {

	w, err := os.Create("list.csv")
	check(err)
	csvW := csv.NewWriter(w)

	err = csvW.Write([]string{"Bezirk", "Stadtteil", "Strasse", "Plaetze", "Wohnart"})
	check(err)

	var wg sync.WaitGroup // async processing of all urls
	wg.Add(len(urls))
	for _, u := range urls {
		go fetchAndParse(&wg, u, csvW)
	}
	wg.Wait()

	csvW.Flush()

	err = csvW.Error()
	check(err)

	err = w.Close()
	check(err)
}

func fetchAndParse(wg *sync.WaitGroup, u string, w *csv.Writer) {
	defer wg.Done()

	doc, err := goquery.NewDocument(u)
	check(err)

	tables := doc.Find("div.richtext table")
	if tables.Length() != 2 {
		check(fmt.Errorf("rce: expected 2 tables on stadtteil-page"))
	}

	t := tables.First()
	if t == nil {
		check(fmt.Errorf("rce: expected non-nil selection"))
	}

	rows := t.Find("tr")

	rows.Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			// skip header row
			return
		}
		fields := row.Find("td")
		err = w.Write([]string{
			fields.Eq(0).Text(),
			fields.Eq(1).Text(),
			fields.Eq(2).Text(),
			fields.Eq(3).Text(),
		})
		check(err)
	})
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
