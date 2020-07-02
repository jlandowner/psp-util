NAME = psp-util
PLUGIN = kubectl-${NAME}
VERSION = v0.1.0
build:
	go build -o bin/${NAME} main.go

build-plugin:
	go build -o bin/${PLUGIN} main.go

update-version:
	sed -i.bak -e "s/v[0-9].[0-9].[0-9]/${VERSION}/" ./cmd/version.go

