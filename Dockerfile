FROM alpine:3.16.3@sha256:b95359c2505145f16c6aa384f9cc74eeff78eb36d308ca4fd902eeeb0a0b161b
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
