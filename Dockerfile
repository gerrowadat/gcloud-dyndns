FROM golang:1.21
RUN go install github.com/gerrowadat/gcloud-dyndns@0.0.3

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

USER root

ENV GCLOUD_DNS_INTERVAL_SECS=86400
ENV GCLOUD_DNS_ZONE=myzone
ENV GCLOUD_DNS_RECORD_NAME=myname.mydomain.tld.
ENV JSON_KEYFILE=/secrets/cloud-dns.key.json

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/entrypoint.sh"]
