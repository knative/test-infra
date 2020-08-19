package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Cluster struct {
	clusterid   string
	accesstoken string
	boskosid    string
	prowid      string
	status      string
}

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=newPassword")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM clusters")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cls := make([]*Cluster, 0)
	for rows.Next() {
		cl := new(Cluster)
		err := rows.Scan(&cl.clusterid, &cl.accesstoken, &cl.boskosid, &cl.prowid, &cl.status)
		if err != nil {
			log.Fatal(err)
		}
		cls = append(cls, cl)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, cl := range cls {
		fmt.Printf("%s, %s, %s, %s, %s\n", cl.clusterid, cl.accesstoken, cl.boskosid, cl.prowid, cl.status)
	}
}
