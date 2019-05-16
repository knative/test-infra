#!/bin/bash

mkdir -p /mount/data/src/
chmod -R a+rw /mount/data

mkdir -p /mount/data/devstats_repos/knative
git clone https://github.com/knative/serving.git /mount/data/devstats_repos/knative/serving

cd /mount/data/src/
git clone https://github.com/ericKlawitter/test-infra.git
cd test-infra/devstats
./scripts/copy_devstats_binaries.sh

rm -rf /etc/gha2db && ln -sf /mount/data/src/test-infra/devstats /etc/gha2db
