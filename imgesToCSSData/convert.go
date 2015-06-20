/* imageToCSSData creates a css file with embedded base64 data from a list of files

$ imageToCSSData my.css Courious_cat.jpg
$ cat my.css

.ipfs-Curious_cat {
	background-image: url('data:image/jpeg;base64,/9j/4AAQSkZJRgABAgAAAQABAA...snip==');
	background-repeat:  no-repeat;
	background-size: contain;
}


*/
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var cssTpl = template.Must(template.New("css").Parse(`
.ipfs-{{.Name}} {
	background-image: url('data:{{.Mime}};base64,{{.Data}}');
	background-repeat:  no-repeat;
	background-size: contain;
}
`))

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: convert [fname.css] [glob of files]")
		os.Exit(1)
	}

	css, err := os.Create(os.Args[1])
	check(err)

	files, err := filepath.Glob(os.Args[2])
	check(err)

	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			raw, err := ioutil.ReadFile(fname)
			check(err)

			shortName := filepath.Base(fname)
			shortName = strings.Split(shortName, ".")[0]

			mtype := http.DetectContentType(raw)

			b64 := base64.StdEncoding.EncodeToString(raw)
			err = cssTpl.Execute(css, struct {
				Name string
				Mime string
				Data string
			}{shortName, mtype, b64})
			check(err)

			fmt.Println(fname, ": converted")
		}(f)
	}

	wg.Wait()
	check(css.Close())
	log.Println("all done")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
