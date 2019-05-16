#!/bin/bash

PGUSER=postgres PGDATABASE=postgres PG_PASS=${PG_PASS} PG_PASS_RO=${PG_PASS} PG_PASS_TEAM=knative ./devel/init_database.sh
./structure

