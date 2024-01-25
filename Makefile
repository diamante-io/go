# Docker build targets use an optional "TAG" environment
# variable can be set to use custom tag name. For example:
#   TAG=my-registry.example.com/keystore:dev make keystore
XDRS = xdr/Diamcircle-SCP.x \
xdr/Diamcircle-ledger-entries.x \
xdr/Diamcircle-ledger.x \
xdr/Diamcircle-overlay.x \
xdr/Diamcircle-transaction.x \
xdr/Diamcircle-types.x

.PHONY: xdr

keystore:
	$(MAKE) -C services/keystore/ docker-build

ticker:
	$(MAKE) -C services/ticker/ docker-build

friendbot:
	$(MAKE) -C services/friendbot/ docker-build

webauth:
	$(MAKE) -C exp/services/webauth/ docker-build

recoverysigner:
	$(MAKE) -C exp/services/recoverysigner/ docker-build

regulated-assets-approval-server:
	$(MAKE) -C services/regulated-assets-approval-server/ docker-build

gxdr/xdr_generated.go: $(XDRS)
	go run github.com/xdrpp/goxdr/cmd/goxdr -p gxdr -enum-comments -o $@ $(XDRS)
	go fmt $@

xdr/xdr_generated.go: $(XDRS) Rakefile Gemfile.lock
	bundle exec rake xdr:generate

xdr: gxdr/xdr_generated.go xdr/xdr_generated.go
