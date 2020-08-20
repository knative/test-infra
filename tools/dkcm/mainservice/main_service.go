package mainservice

import (
	"log"
	"net/http"
	"time"
)

type Options struct {
	Region        string
	BackupRegions []string
	XXTimeout     time.Time
}

type requestInfo struct {
	minNodes int
	maxNodes int
	nodeType string
	zone     string
}

func releaseProject() bool {
	return true
}

func pollProject() string {
	return "Boskos Project"
}

func updateClerk(config requestInfo, prowid string) string {
	return "response token"
}

func handleProw(w http.ResponseWriter, req *http.Request) {
	// handle prow requests

}

func Start(o *Options) {
	log.Printf("Start running the main service with options: %v", o)

	// clustercli, _ := cluster.NewClient()
	// clustercli.Boskos.Release("", "")
	// clerkcli := clerk.NewClient()
}

// func main() {
//
// 	http.HandleFunc("/createcluster", handleProw)
// 	http.HandleFunc("/getcluster", handleProw)
// 	http.ListenAndServe(":8090", nil)
// }
