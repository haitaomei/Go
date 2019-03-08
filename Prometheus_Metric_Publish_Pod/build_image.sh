#!/bin/bash

export dockerUser=
export dockerPwd=

env GOOS=linux GOARCH=amd64 go build

docker login -u $dockerUser -p $dockerPwd

docker build -t khitaomei/custommetric:latest .
docker push khitaomei/custommetric:latest

rm -rf prometheusMetric
