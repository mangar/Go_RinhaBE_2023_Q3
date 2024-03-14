#!/usr/bin/bash

cp ./dockerfile ../../rinha23

cd ../../rinha23

docker build . -t mangar/rinhabe_2023_q3_go:0.0.1

rm dockerfile
