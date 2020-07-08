NAME = psp-util
VERSION = v1.0.1
build:
	go build -o bin/${NAME} main.go

multi-build:
	mkdir -p bin/darwin_amd64 bin/linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64/${NAME} main.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/${NAME} main.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64/${NAME} main.go

update-version:
	sed -i.bak -e "s/v[0-9].[0-9].[0-9][-alpha]*[-beta]*/${VERSION}/g" ./cmd/version.go
	sed -i.bak -e "s/v[0-9].[0-9].[0-9][-alpha]*[-beta]*/${VERSION}/g" ./psp-util.yaml

krew-template:
	docker run -v `pwd`/.krew.yaml:/tmp/template-file.yaml rajatjindal/krew-release-bot:v0.0.38 \
	  krew-release-bot template --tag v1.0.1 --template-file /tmp/template-file.yaml

test-build:
	docker build . -f ./test/test.Dockerfile -t jlandowner/psp-util:${VERSION}
	docker run --rm jlandowner/psp-util:${VERSION}
