package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-mydns"
)

func main() {
	if len(os.Args) != 5 {
		os.Exit(2)
	}

	masterid := os.Getenv("MYDNS_MASTERID")
	password := os.Getenv("MYDNS_PASSWORD")
	if masterid == "" || password == "" {
		fmt.Fprintln(os.Stderr, "$MYDNS_MASTERID and $MYDNS_PASSWORD both must be set")
		os.Exit(2)
	}

	client := mydns.NewClient()
	err := client.Login(masterid, password)
	if err != nil {
		log.Fatal(err)
	}
	di, err := client.FetchRecords()
	if err != nil {
		log.Fatal(err)
	}

	hostname := strings.TrimSuffix(os.Args[2], "."+di.Domain+".")
	log.Printf("Hostname is %v", hostname)
	if os.Args[1] == "present" {
		found := false
		for i := 0; i < len(di.Record); i++ {
			if di.Record[i].Type == "TXT" && di.Record[i].Hostname == hostname {
				log.Printf("updating %v", hostname)
				di.Record[i].Content = os.Args[3]
				found = true
				break
			}
		}
		if !found {
			log.Printf("adding %v", hostname)
			for i := 0; i < len(di.Record); i++ {
				if di.Record[i].Hostname == "" {
					di.Record[i].Hostname = hostname
					di.Record[i].Type = "TXT"
					di.Record[i].Content = os.Args[3]
					break
				}
			}
		}
	} else if os.Args[1] == "cleanup" {
		for i := 0; i < len(di.Record); i++ {
			if di.Record[i].Type == "TXT" && di.Record[i].Hostname == hostname {
				log.Printf("cleanup %v", hostname)
				di.Record[i].Hostname = ""
				di.Record[i].Content = ""
				break
			}
		}
	}

	err = client.UpdateRecords(di)
	if err != nil {
		log.Fatal(err)
	}
}
