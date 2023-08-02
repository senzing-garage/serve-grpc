#!/bin/zsh

export DYLD_LIBRARY_PATH=$LD_LIBRARY_PATH
echo ">>>>>>>>>>>>>>>>>>>> ${LD_LIBRARY_PATH}"
echo ">>>>>>>>>>>>>>>>>>>> ${DYLD_LIBRARY_PATH}"
env


# SENZING_DIR ?= /opt/senzing/g2
# SENZING_TOOLS_SENZING_DIRECTORY ?= $(SENZING_DIR)

# LD_LIBRARY_PATH := $(SENZING_TOOLS_SENZING_DIRECTORY)/lib:$(SENZING_TOOLS_SENZING_DIRECTORY)/lib/macos
# DYLD_LIBRARY_PATH := $(LD_LIBRARY_PATH)

export DYLD_LIBRARY_PATH=/opt/senzing/g2/lib:/opt/senzing/g2/lib/macos:

"$@"