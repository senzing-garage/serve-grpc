#!/bin/zsh

export SENZING_PATH=${HOME}/senzing
export DYLD_LIBRARY_PATH=${SENZING_PATH}/er/lib:${SENZING_PATH}/er/lib/macos
export LD_LIBRARY_PATH=${DYLD_LIBRARY_PATH}
export CGO_CFLAGS="-g -I${SENZING_PATH}/er/sdk/c"
export CGO_LDFLAGS="-L${SENZING_PATH}/er/lib -lSz -Wl,-no_warn_duplicate_libraries"
export SENZING_TOOLS_DATABASE_URL="sqlite3://na:na@nowhere/tmp/sqlite/G2C.db"

"$@"
