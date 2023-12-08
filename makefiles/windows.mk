# Makefile extensions for windows.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: clean-osarch-specific
clean-osarch-specific:
	del /F /S /Q $(TARGET_DIRECTORY)
	del /F /S /Q $(GOPATH)/bin/$(PROGRAM_NAME)
	del /F /S /Q C:\Temp\sqlite


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
	@rmdir /S /Q C:\Temp\sqlite || echo ...but it is not an error.
	@mkdir C:\Temp\sqlite
	@type nul > C:\Temp\sqlite\G2C.db

.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -v -p 1 ./...

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-windows
only-windows:
	@echo "Only windows has this Makefile target."
