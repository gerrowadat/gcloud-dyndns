FROM golang:1.21
RUN go install github.com/gerrowadat/gcloud-dyndns@0.0.2

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

USER root

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/entrypoint.sh"]
