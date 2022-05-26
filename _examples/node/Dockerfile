FROM node:18-alpine

COPY --from=ghcr.io/shopify/hansel:latest /usr/bin/hansel /usr/bin/hansel
RUN hansel --name node --version "$(node -v | sed -e's/^v//g')" --install
