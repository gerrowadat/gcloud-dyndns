#!/usr/bin/bash

set -e

while true
do
  echo "Updating $GCLOUD_DNS_RECORD_NAME in gcloud zone $GCLOUD_DNS_ZONE"
  gcloud-dyndns \
    --cloud-dns-zone=$GCLOUD_DNS_ZONE \
    --cloud-dns-record-name=$GCLOUD_DNS_RECORD_NAME \
    --json-keyfile=$JSON_KEYFILE
  echo "Sleeping for $GCLOUD_DNS_INTERVAL_SECS seconds..."
  sleep $GCLOUD_DNS_INTERVAL_SECS
done
