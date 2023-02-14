FROM alpine:3.17.2@sha256:69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
