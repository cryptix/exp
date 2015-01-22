package main

import (
	"log"
	"os"

	"github.com/cryptix/exp/manageSieve"
)

func main() {
	mc, err := manageSieve.Dial(os.Getenv("MNGSV_HOST"))
	check(err)

	log.Println("Dialed done")

	check(mc.StartTLS(nil))

	log.Println("StartTLS done")

	err = mc.Login(
		os.Getenv("MNGSV_USER"),
		os.Getenv("MNGSV_PASS"))
	check(err)

	scripts, err := mc.ListScripts()
	check(err)
	log.Printf("Scripts: %+v", scripts)

	if len(os.Args) > 1 {
		script, err := mc.GetScript(os.Args[1])
		check(err)
		log.Printf("Script: %s\n%s", os.Args[1], script)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
