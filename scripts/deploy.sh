#!/bin/bash

cd ..
git pull
docker-compose build app
docker-compose up -d