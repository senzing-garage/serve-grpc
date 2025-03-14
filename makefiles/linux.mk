# Makefile extensions for linux.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

LD_LIBRARY_PATH ?= /opt/senzing/er/lib
SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@nowhere/tmp/sqlite/G2C.db
PATH := $(MAKEFILE_DIRECTORY)/bin:/$(HOME)/go/bin:$(PATH)

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: clean-osarch-specific
clean-osarch-specific:
	@docker rm  --force $(DOCKER_CONTAINER_NAME) 2> /dev/null || true
	@docker rmi --force $(DOCKER_IMAGE_NAME) $(DOCKER_BUILD_IMAGE_NAME) $(DOCKER_SUT_IMAGE_NAME) 2> /dev/null || true
	@rm -f  $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm -f  $(MAKEFILE_DIRECTORY)/.coverage || true
	@rm -f  $(MAKEFILE_DIRECTORY)/coverage.html || true
	@rm -f  $(MAKEFILE_DIRECTORY)/coverage.out || true
	@rm -f  $(MAKEFILE_DIRECTORY)/cover.out || true
	@rm -fr $(TARGET_DIRECTORY) || true
	@rm -fr /tmp/sqlite || true
	@pkill godoc || true
	@docker-compose -f docker-compose.test.yaml down 2> /dev/null || true


.PHONY: coverage-osarch-specific
coverage-osarch-specific: export SENZING_LOG_LEVEL=TRACE
coverage-osarch-specific:
	@go test -v -coverprofile=coverage.out -p 1 ./...
	@go tool cover -html="coverage.out" -o coverage.html
	@xdg-open $(MAKEFILE_DIRECTORY)/coverage.html


.PHONY: dependencies-for-development-osarch-specific
dependencies-for-development-osarch-specific:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest


.PHONY: documentation-osarch-specific
documentation-osarch-specific:
	@pkill godoc || true
	@godoc &
	@xdg-open http://localhost:6060


.PHONY: generate-certificates-osarch-specific
generate-certificates-osarch-specific:
	cd testdata/certificates; ./generate_certificates.sh; cd $(MAKEFILE_DIRECTORY)


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	$(info Hello World, from linux.)


.PHONY: package-osarch-specific
package-osarch-specific: docker-build-package
	@mkdir -p $(TARGET_DIRECTORY) || true
	@CONTAINER_ID=$$(docker create $(DOCKER_BUILD_IMAGE_NAME)); \
	docker cp $$CONTAINER_ID:/output/. $(TARGET_DIRECTORY)/; \
	docker rm -v $$CONTAINER_ID


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run -tags "libsqlite3 linux" main.go --enable-all


.PHONY: run-mutual-tls-osarch-specific
run-mutual-tls-osarch-specific: export SENZING_TOOLS_SERVER_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/server/certificate.pem
run-mutual-tls-osarch-specific: export SENZING_TOOLS_SERVER_KEY_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/server/private_key.pem
run-mutual-tls-osarch-specific: export SENZING_TOOLS_CLIENT_CA_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/certificate-authority/certificate.pem
run-mutual-tls-osarch-specific:
	@go run -tags "libsqlite3 linux" main.go --enable-all


.PHONY: run-server-side-tls-osarch-specific
run-server-side-tls-osarch-specific: export SENZING_TOOLS_SERVER_CERTIFICATE_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/server/certificate.pem
run-server-side-tls-osarch-specific: export SENZING_TOOLS_SERVER_KEY_PATH=$(MAKEFILE_DIRECTORY)/testdata/certificates/server/private_key.pem
run-server-side-tls-osarch-specific:
	@go run -tags "libsqlite3 linux" main.go --enable-all


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@mkdir /tmp/sqlite
	@cp testdata/sqlite/G2C.db /tmp/sqlite/G2C.db
	@mkdir -p $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH) || true


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -tags "libsqlite3 linux" -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-linux
only-linux:
	$(info Only linux has this Makefile target.)


.PHONY: run-serve-grpc
run-serve-grpc: build
	@target/linux-amd64/serve-grpc


.PHONY: run-serve-grpc-trace
run-serve-grpc-trace: build
	@target/linux-amd64/serve-grpc --log-level TRACE --engine-log-level 1
