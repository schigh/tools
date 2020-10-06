.PHONY: cuid
cuid:
	go build -o "${GOBIN}/cuid" cli/cmd/cuid/main.go

.PHONY: slug
slug:
	go build -o "${GOBIN}/slug" cli/cmd/slug/main.go
