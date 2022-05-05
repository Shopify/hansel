# Hansel

Hansel generates empty linux packages. These packages can be installed to track dependencies manually added to a container image.

You can use hansel in a multistep build:
```dockerfile
FROM ghcr.io/thepwagner/hansel:latest AS crumbs
RUN hansel --name rando-thing --version v1.2.3 --debian

FROM debian:bullseye
RUN curl -o /usr/bin/rando-thing https://rando.thing/v1.2.3/unsigned-blob-yolo
COPY --from=crumbs /rando-thing*.deb /tmp/rando-thing.deb
RUN dpkg -i /tmp/rando-thing.deb && \
    rm /tmp/rando-thing.deb
```

The name is inspired by:
* [Hansel and Gretel](https://en.wikipedia.org/wiki/Hansel_and_Gretel), as the packages are breadcrumbs left for container scanners to identify.
* [Owen Wilson's character in the 2001 movie "Zoolander"](https://www.youtube.com/watch?v=FAxJECJJG6w), as supply chain observability is "so hot right now".
