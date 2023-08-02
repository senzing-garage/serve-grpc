ifeq ($(OS),Windows_NT)
    OSTYPE = windows
    ifeq ($(PROCESSOR_ARCHITECTURE), AMD64)
        OSARCH = x86_64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE), ARM64)
        OSARCH = arm4
    endif
else
    UNAME_OSTYPE = $(shell uname -s)
    ifeq ($(UNAME_OSTYPE),Linux)
        OSTYPE = linux
        UNAME_ARCH = $(shell uname -p)
        ifeq ($(UNAME_ARCH),x86_64)
            OSARCH = x86_64
        endif
        ifeq ($(UNAME_ARCH),arm64)
            OSARCH = arm64
        endif
    endif
    ifeq ($(UNAME_OSTYPE),Darwin)
        OSTYPE = darwin
        UNAME_ARCH = $(shell uname -m)
        ifeq ($(UNAME_ARCH),x86_64)
            OSARCH = x86_64
        endif
        ifeq ($(UNAME_ARCH),arm64)
            OSARCH = arm64
        endif
    endif
endif
