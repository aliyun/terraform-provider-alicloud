#!/bin/bash

STR=$(go version)

golang1dot8='go1.8'
golang1dot9='go1.9'
golang1dot10='go1.10'
golang1dot11='go1.11'

function setup_golang_vendor() {
  # Use glide to install deps for go version < 1.11
  curl https://glide.sh/get | sh
  glide update
  glide install
  go test -v
}

function setup_golang_go_mod() {
  # Use go mod for go version >= 1.11
  export GO111MODULE=on
  go test -v
}

if [[ "$STR" =~ .*"$golang1dot8".* ]]; then
  echo "Setup golang 1.8"
  setup_golang_vendor
fi

if [[ "$STR" =~ .*"$golang1dot9".* ]]; then
  echo "Setup golang 1.9"
  setup_golang_vendor
fi

if [[ "$STR" =~ .*"$golang1dot10".* ]]; then
  echo "Setup golang 1.10"
  setup_golang_vendor
fi

if [[ "$STR" =~ .*"$golang1dot11".* ]]; then
  echo "Setup golang 1.11"
  setup_golang_go_mod
fi

