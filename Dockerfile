FROM alpine:3.20.3@sha256:1e42bbe2508154c9126d48c2b8a75420c3544343bf86fd041fb7527e017a4b4a
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
