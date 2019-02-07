#!/usr/bin/env bash

go test -coverprofile new_profile.txt
$2 download knative-prow $1 -o base_profile.txt
$2 diff base_profile.txt new_profile.txt