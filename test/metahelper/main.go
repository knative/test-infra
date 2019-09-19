package main

import (
	"flag"
	"log"

	"knative.dev/test-infra/test/metahelper/util"
)

func main() string {
	getKeyOpt := flag.String("get", "", "get val for a key")
	saveKeyOpt := flag.String("set", "", "save val for a key, must have --val supplied")
	valOpt := flag.String("val", "", "val to be modified, only useful when --save is passed")
	flag.Parse()
	// Create with default path
	c, err := util.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *getKeyOpt != "" && *saveKeyOpt != "":
		log.Fatal("--get and --save can't be used at the same time")
	case *getKeyOpt != "":
		gotVal, err := c.Get(*getKeyOpt)
		if err != nil {
			log.Fatalf(err)
		}
		return gotVal
	case *saveKeyOpt != "":
		if *valOpt == "" {
			log.Fatal("--val must be supplied when using --save")
		}
		log.Printf("Writing files to %s", c.Path)
		if err := c.Set(*saveKeyOpt, *valOpt); err != nil {
			log.Fatalf(err)
		}
		return ""
	}
}
