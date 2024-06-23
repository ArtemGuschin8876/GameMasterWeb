#!/bin/sh
migrate -path migrations/ -database "postgresql://postgres:postgres@postgres:5432/gamemasterweb?sslmode=disable" -verbose up && ./api
