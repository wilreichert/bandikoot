#!/bin/bash -uex

cd python
docker build -t python:2.7-alpine3.5 .
cd ..
