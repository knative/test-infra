package mainservice

import "time"

type clusterInfo struct {
	clusterName string
	boskosid    string
	zone        string
}

func createCluster(clusterName string, boskosid string, zone string, update chan bool) *clusterInfo {
	// if in database, check status
	time.Sleep(30 * time.Second)
	update <- true
	return &clusterInfo{clusterName, boskosid, zone}
}

func clearCluster(clusterName string, boskosid string, zone string, update chan bool) {
	time.Sleep(3 * time.Second)
	update <- true
}
