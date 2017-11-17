#!/bin/bash

pushd /opt/creative_info_manager

killall creative_info_manager

sleep 1

nohup bin/creative_info_manager > /pdata/logs/creative_info_manager/creative.log 2>&1 &

sleep 1

popd

