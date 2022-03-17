.PHONY: cuid
cuid:
	go build -o "${GOBIN}/cuid" cli/cmd/cuid/main.go

.PHONY: slug
slug:
	go build -o "${GOBIN}/slug" cli/cmd/slug/main.go

.PHONY: uuid
uuid:
	go build -o "${GOBIN}/uuid" cli/cmd/uuid/main.go

.PHONY: uuidv1
uuidv1:
	go build -o "${GOBIN}/uuidv1" cli/cmd/uuidv1/main.go
.PHONY: md5
md5:
	go build -o "${GOBIN}/md5" cli/cmd/md5/main.go

.PHONY: sha1
sha1:
	go build -o "${GOBIN}/sha1" cli/cmd/sha1/main.go

.PHONY: sha256
sha256:
	go build -o "${GOBIN}/sha256" cli/cmd/sha256/main.go

.PHONY: bcrypt
bcrypt:
	go build -o "${GOBIN}/bcrypt" cli/cmd/bcrypt/main.go

.PHONY: guid
guid:
	go build -o "${GOBIN}/guid" cli/cmd/guid/main.go

.PHONY: all
all: cuid slug uuid uuidv1 md5 sha1 sha256 bcrypt guid
