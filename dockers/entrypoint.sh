#!/bin/sh

set -e

# envoy will listen on port 8080
nohup server -address 0.0.0.0 &

/usr/local/bin/envoy -l debug -c "/envoy.yaml"