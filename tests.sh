#!/usr/bin/env sh

function run_tests() {
  echo "Running tests for $1"
  make test-integration WHAT=./test/integration/$1 GOFLAGS="-v -failfast -run=$2" &> $1.out
}

run_tests statefulset
run_tests deployment
run_tests replicaset
run_tests scheduler "TestUnschedulableNodes|TestAllocatable"
run_tests job
