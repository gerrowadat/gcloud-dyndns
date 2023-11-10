#!/usr/bin/bash

set -e

GCLOUD_DNS_INTERVAL_SECS=86400

while true
do
  gcloud-dyndns "$*"
  sleep $CLOUD_DNS_INTERVAL_SECS
done
