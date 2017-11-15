#!/bin/bash

if [ ! -d "/pdata/logs/creative_info_manager" ]; then
    mkdir -p /pdata/logs/creative_info_manager
fi

pushd /opt/creative_info_manager

if [ ! -d "logs" ]; then
    mkdir logs
fi

make deps > build.log 2>&1 || (cat build.log && exit 1)
make > build.log 2>&1 || (cat build.log && exit 1)
    
popd
