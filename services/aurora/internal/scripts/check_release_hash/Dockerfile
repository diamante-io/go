# Change to Go version used in CI or rebuild with --build-arg.
ARG GO_IMAGE=golang:1.17
FROM $GO_IMAGE

WORKDIR /go/src/github.com/diamcircle/go

ENV DEBIAN_FRONTEND=noninteractive
# ca-certificates are required to make tls connections
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl wget gnupg apt-utils git zip unzip apt-transport-https ca-certificates
RUN wget -qO - https://apt.diamcircle.org/SDF.asc | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=true apt-key add -
RUN echo "deb https://apt.diamcircle.org xenial stable" >/etc/apt/sources.list.d/SDF.list
RUN echo "deb https://apt.diamcircle.org xenial testing" >/etc/apt/sources.list.d/SDF-testing.list

RUN git clone https://github.com/diamcircle/go.git /go/src/github.com/diamcircle/go
# Fetch dependencies and prebuild binaries. Not necessary but will make check faster.
RUN go run -v ./support/scripts/build_release_artifacts

COPY check.sh .
RUN chmod +x check.sh
ENTRYPOINT ["./check.sh"]
