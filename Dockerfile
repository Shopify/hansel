FROM alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
