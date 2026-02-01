statictests:
	go vet -vettool=statictest ./...

build:
	go build -ldflags "-X main.Version=v1.0.1 -X main.BuildCommit=$(git rev-parse --short HEAD) -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/client/main.go