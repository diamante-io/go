SUDO := $(shell docker version >/dev/null 2>&1 || echo "sudo")

# https://github.com/opencontainers/image-spec/blob/master/annotations.md
BUILD_DATE := $(shell date -u +%FT%TZ)

TAG ?= diamcircle/diamcircle-aurora:$(VERSION)

docker-build:
ifndef VERSION
	$(error VERSION environment variable must be set. For example VERSION=2.4.1-101 )
endif
	$(SUDO) docker build --pull \
	--label org.opencontainers.image.created="$(BUILD_DATE)" \
	--build-arg VERSION=$(VERSION) \
	-t $(TAG) .

docker-push:
ifndef TAG
	$(error Must set VERSION or TAG environment variable. For example VERSION=2.4.1-101 )
endif
	$(SUDO) docker push $(TAG)
