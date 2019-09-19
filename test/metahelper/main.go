package main

import (
	"flag"
	"log"

	"knative.dev/eventing-operator/test/metahelper/util"
)

func main() {
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
			log.Print("")
		} else {
			log.Print(gotVal)
		}
	case *saveKeyOpt != "":
		if *valOpt == "" {
			log.Fatal("--val must be supplied when using --save")
		}
		err := c.Set(*saveKeyOpt, *valOpt)
		if err != nil {
			log.Print(err)
		} else {
			log.Print("")
		}
	}

}
