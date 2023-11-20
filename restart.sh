#!/bin/bash

make stop
sleep 2
make start
sleep 6
make fix-start
sleep 5
make bff-container-logs
