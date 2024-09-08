# needs to be run in a go environment

run:
	go run cmd/main/main.go $(ARGS)

build:
	go run cmd/build/build.go

release:
	go run cmd/release/release.go

get_version:
	go run cmd/version/version.go
