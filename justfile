alias r := run
alias b := build
alias t := test

_default:
  just -l

run:
  go run pkg/main.go

build:
  go build -o main ./pkg

test pkg:
	fd {{pkg}} | xargs -I# go test ./#
