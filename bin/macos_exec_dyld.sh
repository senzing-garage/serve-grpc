#!/bin/zsh

export DYLD_LIBRARY_PATH=$LD_LIBRARY_PATH
echo ">>>>>>>>>>>>>>>>>>>> ${LD_LIBRARY_PATH}"
echo ">>>>>>>>>>>>>>>>>>>> ${DYLD_LIBRARY_PATH}"
env

"$@"