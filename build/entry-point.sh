#!/bin/sh

set -e

case "$1" in
    "spider")
        echo "Running fly-go spider..."
        exec /usr/local/bin/fly-go spider
        ;;
    "server")
        echo "Running fly-go server..."
        exec /usr/local/bin/fly-go server
        ;;
    "all")
        echo "Running fly-go all..."
        exec /usr/local/bin/fly-go
        ;;
    *)
        echo "Running default sh..."
        exec sh
        ;;
esac