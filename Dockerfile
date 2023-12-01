FROM alpine:3.18.5@sha256:34871e7290500828b39e22294660bee86d966bc0017544e848dd9a255cdf59e0
ENTRYPOINT ["/usr/bin/hansel"]
COPY hansel /usr/bin/hansel
