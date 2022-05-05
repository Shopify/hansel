FROM alpine:3.15.4
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
