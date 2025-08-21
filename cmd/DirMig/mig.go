package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/k3mpton/shortner-project/pkg/conndb"
	"github.com/pressly/goose"
)

var (
	DirectitonMig = flag.String("mig", "up", "direction migration 'up' or 'down'")
)

func main() {
	flag.Parse()
	mig := strings.ToLower(*DirectitonMig)

	if mig != "up" && mig != "down" {
		log.Fatalf("failed direct migration, fail argument: %v", mig)
	}

	db := conndb.Conn()

	direct := "./migrations"
	fmt.Println("hrllo")
	switch mig {
	case "up":
		if err := goose.Up(db, direct); err != nil {
			log.Fatalf("fail, migration to up: %v", err)
		}
	default:
		if err := goose.Down(db, direct); err != nil {
			log.Fatalf("fail, migration to down: %v", err)
		}
	}
}
