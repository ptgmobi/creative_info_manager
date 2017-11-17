#!/bin/bash

if [ ! -d "/pdata/logs/creative_info_manager" ]; then
    mkdir -p /pdata/logs/creative_info_manager
fi

pushd /opt/creative_info_manager

aws s3 sync s3://cloudmobi-config/creative_info_manager/conf conf/ 

if [ ! -d "logs" ]; then
    mkdir logs
fi

make deps > build.log 2>&1 || (cat build.log && exit 1)
make > build.log 2>&1 || (cat build.log && exit 1)
    
popd
