#!/usr/bin/env bash

source "${HOME}/.gvm/scripts/gvm"

if [[ -v GO_VERSION ]]; then
  gvm use "${GO_VERSION}"
fi

kubekins-runner.sh "$@"
