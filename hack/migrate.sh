#!/usr/bin/env bash

do_migrate() {
    migrate create -ext sql -dir "database/$1/migrations" "$2"
}

usage() {
    echo "Usage:"
    exit 1
}

case "$1" in
    postgres)
        do_migrate postgres "$2"
        ;;
    sqlite)
        do_migrate sqlite "$2"
        ;;
    mysql)
        do_migrate mysql "$2"
        ;;
    mariadb)
        do_migrate mariadb "$2"
        ;;
    all)
        for i in mariadb postgres sqlite mysql; do
            do_migrate "$i" "$2"
        done
        ;;
    *)
        usage
        ;;
esac
