#!/usr/bin/env sh

migrate create -ext sql -dir "$1/migrations" "$2"
