# needs to be run in a go environment

build:
	go run build-scripts/build/build.go

release:
	go run build-scripts/release/release.go

get_version:
	go run build-scripts/version/version.go