FROM alpine:3.16.2@sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
