/*
table generator for

	https://wiki.freifunk.net/index.php?title=Hamburg/Fl%C3%BCchtlinge#Liste_von_Unterk.C3.BCnften

TODO

1. consolidate existing data (status, doku links, contact person)

2. intermediate datastorage to detect upstream changes

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// from: https://www.hamburg.de/fluechtlinge-unterbringung-standorte/
var urls = []string{
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4372596/unterbringung-altona/",
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373118/unterbringung-bergedorf/",
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373120/unterbringung-eimsbuettel/",
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373122/unterbringung-harburg/",
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373128/unterbringung-mitte/",
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373126/unterbringung-nord/",
	"https://www.hamburg.de/fluechtlinge-unterbringung-standorte/4373124/unterbringung-wandsbek/",
}

var re = regexp.MustCompile(`\(Stand ([0-9]+\.[0-9]+\.[0-9]+)\)`)

type result struct {
	Bezirk, Stadtteil string
	Strasse, Plaetze  string
	Wohnart           string

	Date time.Time

	err error
}

var csvFname = flag.String("fname", "list.tbl", "filename to write to")

var mediaWikiTpl = template.Must(template.New("mediaTable").Parse(`
{| class="mw-datatable sortable toptextcells"
! '''Bezirk'''
! '''Stadtteil'''
! '''Straße'''
! '''Plätze'''
! '''Wohnart'''
! '''Anprechpartner*'''
! '''Doku'''
! '''Status**'''
{{range .}}
|-
| {{.Bezirk}}
| {{.Stadtteil}}
| {{.Strasse}}
| {{.Plaetze}}
| {{.Wohnart}}
| TODOKontakt
| TODODoku
| TODOStatus{{end}}
|}
`))

func main() {
	flag.Parse()

	w, err := os.Create(*csvFname)
	check(err)

	var rs []result
	for r := range fetchAndParse(urls) {
		check(r.err)
		rs = append(rs, r)
	}

	err = mediaWikiTpl.Execute(w, rs)
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
				r.Bezirk = strings.Title(u[i+len(p) : len(u)-1])

				introText := doc.Find("p.intro").Text()
				m := re.FindStringSubmatch(introText)
				if len(m) != 2 {
					r.err = fmt.Errorf("rce: did not find (Stand $date) in html document")
					c <- r
					return
				}
				r.Date, err = time.Parse("02.01.2006", m[1])
				if err != nil {
					r.err = err
					c <- r
					return
				}

				tables := doc.Find("div.richtext table")
				if tables.Length() != 2 {
					r.err = fmt.Errorf("rce: expected 2 tables on Stadtteil-page")
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

				log.Printf("%s: Bestehende Standorte: %d (Stand: %s)", r.Bezirk, rows.Length()-1, r.Date.Format("02.01.2006"))

				rows.Each(func(i int, row *goquery.Selection) {
					if i == 0 {
						// skip header row
						return
					}
					fields := row.Find("td")
					r.Stadtteil = fields.Eq(0).Text()
					r.Strasse = fields.Eq(1).Text()
					r.Plaetze = fields.Eq(2).Text()
					r.Wohnart = fields.Eq(3).Text()
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
