#!/usr/bin/env bash

./scripts/build.sh && serverless deploy --stage $1 
