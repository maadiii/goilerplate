# goilerplate
A clean architecture and structure for golang web projects

## config database
This is use sql database (postgres, mysql, sqlserver,.... etc).

create user with goilerplate name and goilerplate password on your database.

create database with goilerplate name and give all privileg to that user already created

of curse you can chang database settings. go to ./db/consts.go and change constant with DEBUG prefix


## run tests
run go mod tidy to commit dependencies

run go test -v ./... -p1 for pass all tests


## run application
type go run main.go serve or go run main.go serve --config ./config.yml for not debug mode.

open browser and go adress http://localhost:8000/admin/groups


## more help
type go run main.go --help for more helps.


I know this tutorial not enough but I'll complete this soon.
