FROM alpine:3.18.2@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
