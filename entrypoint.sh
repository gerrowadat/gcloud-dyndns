#!/usr/bin/bash

set -e

CLOUD_DNS_INTERVAL_SECS=86400

GCLOUD_PROJECT=myproject
GCLOUD_DNS_ZONE=myzone
GCLOUD_DNS_RECORD_NAME=myname.mydomain.tld.
JSON_KEYFILE=/secrets/cloud-dns.key.json


while true
do
  gcloud-dyndns --cloud-project=$GCLOUD_PROJECT \
    --cloud-dns-zone=$GCLOUD_DNS_ZONE \
    --cloud-dns-record-name=$GCLOUD_DNS_RECORD_NAME \
    --json-keyfile=$JSON_KEYFILE
  sleep $CLOUD_DNS_INTERVAL_SECS
done
