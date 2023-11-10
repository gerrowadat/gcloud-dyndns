FROM golang:1.21
RUN go install github.com/gerrowadat/gcloud-dyndns@54f9b81

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

USER root

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/entrypoint.sh"]
