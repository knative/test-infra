package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
	data := url.Values{}
	resp, err := http.PostForm("http://localhost:8080/acquire?type=prow-cluster&state=free&dest=busy&owner=prow", data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res)
	return fmt.Sprint(res["name"])
}

func updateClerk(config requestInfo, prowid string) string {
	tokenChannel := make(chan string)
	boskosid := pollProject()
	go updateCluster(config.zone, prowid, boskosid, tokenChannel)
	return <-tokenChannel
}

func handleProw(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.Write([]byte("we are checking whether your cluster has been created\n"))
		var token string
		for k, v := range req.URL.Query() {
			fmt.Printf("%s: %s\n", k, v)
			if k == "token" {
				token = v[0]
			}
		}
		statusChannel := make(chan string, 5)
		errorChannel := make(chan string, 1)
		go getWithToken(token, statusChannel, errorChannel)
	ForLoop:
		for {
			select {
			case err := <-errorChannel:
				w.Write([]byte(err))
				break ForLoop
			case status := <-statusChannel:
				w.Write([]byte(status))
				for val := range statusChannel {
					w.Write([]byte(val))
				}
				break ForLoop
			default:
				w.Write([]byte("waiting for clerk to verify...\n"))
				time.Sleep(1 * time.Millisecond)
			}
		}
	case "POST":
		prowid := req.PostFormValue("prowid")
		minNodes, err1 := strconv.Atoi(req.PostFormValue("minNodes"))
		maxNodes, err2 := strconv.Atoi(req.PostFormValue("maxNodes"))
		if err1 != nil {
			panic(err1)
		}
		if err2 != nil {
			panic(err2)
		}
		nodeType := req.PostFormValue("nodeType")
		zone := req.PostFormValue("zone")
		configInfo := requestInfo{minNodes, maxNodes, nodeType, zone}
		received := "configuration: nodeType: " + nodeType + " prow id:" + prowid + " zone: " + zone + "\n"
		//fmt.Println(received)
		w.Write([]byte(received))
		accessToken := updateClerk(configInfo, prowid)
		w.Write([]byte("your accesssToken to the cluster is : " + accessToken + "\n"))
	}
	clerkInfo()

}

func main() {
	http.HandleFunc("/createcluster", handleProw)
	http.HandleFunc("/getcluster", handleProw)
	http.ListenAndServe(":8090", nil)
}
