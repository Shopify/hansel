FROM debian:buster-slim AS debian
COPY tmp/*.deb /test.deb
RUN dpkg -i /test.deb && \
  dpkg -l | grep -q "ii  meow"

FROM alpine:3.15.4 AS alpine
COPY tmp/*.apk /test.apk
RUN apk add --allow-untrusted --no-network /test.apk
