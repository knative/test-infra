package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Cluster struct {
	clusterid   string
	accesstoken string
	boskosid    string
	prowid      string
	status      string
	zone        string
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	db        *sql.DB
	randomGen *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=newPassword")
	if err != nil {
		log.Fatal(err)
	}
	// check connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

}

func generateID(idSize int) string {
	bytes := make([]byte, idSize)
	for i := range bytes {
		bytes[i] = charset[randomGen.Intn(len(charset))]
	}
	return string(bytes)
}

func generateUnique(idSize int, key string) string {
	var randomid string
	var count int
	for {
		randomid = generateID(idSize)
		db.QueryRow("SELECT count(*) FROM table Where $1 = $2", key, randomid).Scan(&count)
		if count == 0 {
			return randomid
		}
	}
}

func clerkInfo() {
	rows, err := db.Query("SELECT * FROM clusters")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cls := make([]*Cluster, 0)
	for rows.Next() {
		cl := new(Cluster)
		err := rows.Scan(&cl.clusterid, &cl.accesstoken, &cl.boskosid, &cl.prowid, &cl.status, &cl.zone)
		if err != nil {
			log.Fatal(err)
		}
		cls = append(cls, cl)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, cl := range cls {
		fmt.Printf("%s, %s, %s, %s, %s, %s\n", cl.clusterid, cl.accesstoken, cl.boskosid, cl.prowid, cl.status, cl.zone)
	}
}

func clerkQuery() (bool, string, string) {
	row := db.QueryRow("SELECT * FROM clusters WHERE status = 'ready' AND prowid = '0' ")
	cl := new(Cluster)
	err := row.Scan(&cl.clusterid, &cl.accesstoken, &cl.boskosid, &cl.prowid, &cl.status, &cl.zone)
	if err == sql.ErrNoRows {
		fmt.Println("no cluster available. We are creating it!\n")
		return false, "", ""
	} else if err != nil {
		fmt.Println(err)
		return false, "", ""
	}
	fmt.Println("cluster found!")
	fmt.Printf("%s, %s, %s, %s, %s\n", cl.clusterid, cl.accesstoken, cl.boskosid, cl.prowid, cl.status, cl.zone)
	return true, cl.clusterid, cl.accesstoken
}

func getWithToken(token string, status chan string, errc chan string) {
	row := db.QueryRow("SELECT * FROM clusters WHERE accesstoken = $1", token)
	cl := new(Cluster)
	err := row.Scan(&cl.clusterid, &cl.accesstoken, &cl.boskosid, &cl.prowid, &cl.status, &cl.zone)
	if err == sql.ErrNoRows {
		errc <- "Illegal token! No trepassing!\n"
	} else if err != nil {
		errc <- "We don't understand your request. Please try again.\n"
	} else {
		if strings.Trim(cl.status, " ") == "ready" {
			status <- "cluster is ready! Thank you for your patience! Here is the info: \n"
			status <- "cluster_name: " + cl.clusterid + "\n"
			status <- "project_id: " + cl.boskosid + "\n"
			status <- "zone: " + cl.zone + "\n"
		} else {
			status <- "creation is still in progress. Patience :) \n"
		}

	}
	close(status)

}

func updateCluster(zone string, prowid string, boskosid string, infoChan chan string) {
	avail, clusterid, accesstoken := clerkQuery()
	if avail == true {
		fmt.Println("We have found available cluster, now organizing info...\n")
		fmt.Println(accesstoken)
		infoChan <- accesstoken
		fmt.Println("hello")
		updateStatement := `UPDATE clusters
		SET prowid = $2
		WHERE clusterid = $1;`
		_, err := db.Exec(updateStatement, clusterid, prowid)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("We have not found available cluster, now creating new...\n")
		createStatement := `INSERT INTO
		clusters VALUES($1, $2, $3, $4, $5, $6);
		`
		var clusterid = generateUnique(6, "clusterid")
		accesstoken = generateUnique(9, "accesstoken")
		infoChan <- accesstoken
		_, err := db.Exec(createStatement, clusterid, accesstoken, boskosid, prowid, "in progress", zone)
		if err != nil {
			panic(err)
		}
		createUpdate := make(chan bool)
		go createCluster(clusterid, boskosid, zone, createUpdate)
		update := <-createUpdate
		if update == true {
			updateStatement := `UPDATE clusters
			SET status = $2
			WHERE clusterid = $1;`
			_, err := db.Exec(updateStatement, clusterid, "ready")
			if err != nil {
				panic(err)
			}
		}
	}

}
