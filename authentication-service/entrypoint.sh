#! /bin/sh

migrate -path=./migrations -database=postgresql://db_admin:pass5@word@postgres/users up
