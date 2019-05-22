#!/bin/bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Examples:
# MODE=rr FROM=`cat input` TO=`cat output` FILES=`find abc/ -type f -not -iname 'something.txt'` ./devel/mass_replace.sh
# MODE=ss FROM=`cat input` TO=`cat output` FILES=`ls grafana/dashboards/{all,cncf,cni,containerd,coredns,envoy,fluentd,grpc,jaeger,linkerd,kubernetes,notary,opencontainers,opentracing,prometheus,rkt,rook,tuf,vitess}/*` ./devel/mass_replace.sh
# MODE=ss0 FROM=CNCF TO='[[full_name]]' FILES=`find grafana/dashboards/cncf/ -type f -not -iname dashboards.json` ./devel/mass_replace.sh
# MODE=ss FROM='      "title": "' TO='      "title": "[[full_name]] ' FILES=`find grafana/dashboards/cncf/ -name "top_commenters.json" -or -name "project_statistics.json" -or -name "companies_summary.json" -or -name "prs_authors_companies_histogram.json" -or -name "developers_summary.json" -or -name "prs_authors_histogram.json"` ./devel/mass_replace.sh
# MODE=rs0 FROM='(?m)^.*"uid": "\w+",\n' TO='-' ./replacer input.json
# MODE=rr0 FROM='(?m)(^.*)"uid": "(\w+)",' TO='$1"uid": "placeholder",' ./replacer input.json
# MODE=rr FROM='(?m);;;(.*)$' TO=';;;$1 # {{repo_groups}}' FILES=`find metrics/ -iname "gaps.yaml"` ./devel/mass_replace.sh
if [ -z "${FROM}" ]
then
  echo "You need to set FROM, example FROM=abc TO=xyz FILES='f1 f2' $0"
  exit 1
fi
if [ -z "${TO}" ]
then
  echo "You need to set TO, example FROM=abc TO=xyz FILES='f1 f2' $0"
  exit 2
fi
if [ -z "${FILES}" ]
then
  echo "You need to set FILES, example FROM=abc TO=xyz FILES='f1 f2' $0"
  exit 3
fi
for f in ${FILES}
do
  ./replacer $f || exit 4
done
echo 'OK'
