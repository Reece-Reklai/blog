#!/usr/bin/env bash

DRIVER=postgres
DBSTRING="postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

cd sql/schema/ || exit 1

if [[ "$1" = "up" || "$1" = "down" ]]; then
    goose "$DRIVER" "$DBSTRING" "$1"
else
    echo "Requires argument 'up' or 'down' for now . . ."
    exit 1
fi
