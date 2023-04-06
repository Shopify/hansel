FROM alpine:3.17.3@sha256:124c7d2707904eea7431fffe91522a01e5a861a624ee31d03372cc1d138a3126
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
