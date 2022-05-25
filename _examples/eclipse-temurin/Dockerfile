FROM eclipse-temurin:18.0.1_10-jre-alpine

# In this example, we add hansel to the image and then automatically install the package in a single step
COPY --from=ghcr.io/shopify/hansel:latest /usr/bin/hansel /usr/bin/hansel
RUN hansel \
  --name openjdk \
  --version $(java -version 2>&1 | grep "Runtime Environment" | cut -d' ' -f4) \
  --install
