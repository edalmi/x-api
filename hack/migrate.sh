#!/usr/bin/env sh

migrate create -ext sql -dir "database/$1/migrations" "$2"
