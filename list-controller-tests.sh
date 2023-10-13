#!/usr/bin/env bash

set -x

rm tests.txt

for controller in scheduler cronjob daemonset deployment garbagecollector job node podgc pods replicaset statefulset; do
  echo "===== test/integration/$controller" >> tests.txt
  go test ./test/integration/$controller -list=. >> tests.txt
done
