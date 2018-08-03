#!/bin/bash


env GOOS=linux GOARCH=amd64 go build

docker login -u khitaomei -p galaxy123

docker build -t khitaomei/custommetric:latest .
docker push khitaomei/custommetric:latest

rm -rf prometheusMetric