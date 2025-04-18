alias r := run
alias b := build

_default:
  just -l

run:
  go run pkg/main.go

build:
  go build -o main ./pkg
