#!/bin/env sh

./go-app -migrate=true
./go-app -seed=true
./go-app -seedsample=true
./go-app
