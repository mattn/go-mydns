package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-mydns"
)

type stringList []string

func (ss *stringList) String() string {
	return "my string representation"
}

func (ss *stringList) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}

func main() {
	var addTXT stringList
	var updTXT stringList
	flag.Var(&addTXT, "at", "add new TXT entry (-at foo=bar)")
	flag.Var(&updTXT, "ut", "update TXT entry (-ut foo=bar)")
	flag.Parse()

	if len(updTXT) == 0 && len(addTXT) == 0 {
		flag.Usage()
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

	for _, arg := range updTXT {
		tok := strings.SplitN(arg, "=", 2)
		if len(tok) == 1 {
			tok = append(tok, "")
		}
		for i := 0; i < len(di.Record); i++ {
			if di.Record[i].Type == "TXT" && di.Record[i].Hostname == tok[0] {
				di.Record[i].Content = tok[1]
			}
		}
	}

	for _, arg := range addTXT {
		tok := strings.SplitN(arg, "=", 2)
		if len(tok) == 1 {
			tok = append(tok, "")
		}
		for i := 0; i < len(di.Record); i++ {
			if di.Record[i].Hostname == "" {
				di.Record[i].Hostname = tok[0]
				di.Record[i].Type = "TXT"
				di.Record[i].Content = tok[1]
			}
		}
	}

	err = client.UpdateRecords(di)
	if err != nil {
		log.Fatal(err)
	}
}
