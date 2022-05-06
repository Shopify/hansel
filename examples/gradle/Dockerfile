# In this example, we generate the metapackages in a separate build step
# Also "hansel and gradle" made me laugh

ARG JDK_VERSION=18.0.1_10
ARG GRADLE_VERSION=7.4.2
ARG GRADLE_DIGEST=47803a7139e779f97eeb87f35f5980eb57040535861f746bf9efeffa7d1dd98ae3b3296d1d180444e21eadcdbeac748dc493ae9a2cff552a5629455bbbaa9074

FROM ghcr.io/Shopify/hansel:latest AS crumbs
ARG JDK_VERSION
ARG GRADLE_VERSION
RUN hansel --name openjdk --version $(echo $JDK_VERSION | sed -e 's/_/-/g') --debian
RUN hansel --name gradle --version $GRADLE_VERSION --debian

FROM eclipse-temurin:${JDK_VERSION}-jre
ARG GRADLE_VERSION
ARG GRADLE_DIGEST
RUN apt-get update && \
    apt-get install -y \
        unzip \
        && \
    rm -rf /var/lib/apt/lists/*
RUN curl -sLo /gradle.zip "https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip" && \
    echo "${GRADLE_DIGEST}  /gradle.zip" | sha512sum -c - && \
    unzip -o /gradle.zip -d /usr/local/gradle && \
    rm /gradle.zip && \
    ln -s /usr/local/gradle/gradle-${GRADLE_VERSION}/bin/gradle /usr/bin/gradle

COPY --from=crumbs /openjdk*.deb /openjdk.deb
COPY --from=crumbs /gradle*.deb /gradle.deb
RUN dpkg -i /openjdk.deb && \
    dpkg -i /gradle.deb && \
    rm /openjdk.deb /gradle.deb
