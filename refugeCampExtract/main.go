package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

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

var re = regexp.MustCompile(`\(Stand ([0-9]+\.[0-9]+\.[0-9]+)\)`)

type result struct {
	bezirk, stadtteil string
	strasse, plaetze  string
	wohnart           string

	date time.Time

	err error
}

var csvFname = flag.String("fname", "list.csv", "filename to write to")

func main() {
	flag.Parse()

	w, err := os.Create(*csvFname)
	check(err)

	csvW := csv.NewWriter(w)

	// csv header row
	err = csvW.Write([]string{"Bezirk", "Stadtteil", "Strasse", "Plaetze", "Wohnart"})
	check(err)

	for r := range fetchAndParse(urls) {
		check(r.err)

		err = csvW.Write([]string{
			r.bezirk,
			r.stadtteil,
			r.strasse,
			r.plaetze,
			r.wohnart,
		})
		check(err)
	}

	csvW.Flush()

	err = csvW.Error()
	check(err)

	err = w.Close()
	check(err)
}

func fetchAndParse(urls []string) <-chan result {
	c := make(chan result)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(urls))

		for _, u := range urls {
			go func(u string) {
				defer wg.Done()
				var r result

				doc, err := goquery.NewDocument(u)
				if err != nil {
					r.err = err
					c <- r
					return
				}

				// strip out the Bezirk from the url
				var p = "unterbringung-"
				i := strings.LastIndex(u, p)
				r.bezirk = strings.Title(u[i+len(p) : len(u)-1])

				introText := doc.Find("p.intro").Text()
				m := re.FindStringSubmatch(introText)
				if len(m) != 2 {
					r.err = fmt.Errorf("rce: did not find (Stand $date) in html document")
					c <- r
					return
				}
				r.date, err = time.Parse("02.01.2006", m[1])
				if err != nil {
					r.err = err
					c <- r
					return
				}

				tables := doc.Find("div.richtext table")
				if tables.Length() != 2 {
					r.err = fmt.Errorf("rce: expected 2 tables on stadtteil-page")
					c <- r
					return
				}

				t := tables.First()
				if t == nil {
					r.err = fmt.Errorf("rce: expected non-nil selection")
					c <- r
					return
				}

				rows := t.Find("tr")

				log.Printf("%s: Bestehende Standorte: %d (Stand: %s)", r.bezirk, rows.Length()-1, r.date.Format("02.01.2006"))

				rows.Each(func(i int, row *goquery.Selection) {
					if i == 0 {
						// skip header row
						return
					}
					fields := row.Find("td")

					r.stadtteil = fields.Eq(0).Text()
					r.strasse = fields.Eq(1).Text()
					r.plaetze = fields.Eq(2).Text()
					r.wohnart = fields.Eq(3).Text()
					c <- r
				})
			}(u)
		}
		wg.Wait()
		close(c)
	}()
	return c
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
