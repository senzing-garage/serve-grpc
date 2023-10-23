# Makefile extensions for windows.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------


# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: windows/amd64
	@mv $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME).exe


.PHONY: clean-osarch-specific
clean-osarch-specific:
	del /F /S /Q $(TARGET_DIRECTORY)
	del /F /S /Q $(GOPATH)/bin/$(PROGRAM_NAME)


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	@echo "Hello World, from windows."


.PHONY: package-osarch-specific
package-osarch-specific:
	@echo No packaging for windows.


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@echo "No setup required."


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -v -p 1 ./...

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-windows
only-windows:
	@echo "Only windows has this Makefile target."
