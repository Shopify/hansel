FROM alpine:3.17.1@sha256:f271e74b17ced29b915d351685fd4644785c6d1559dd1f2d4189a5e851ef753a
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
