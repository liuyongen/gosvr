BINARY = Pchat-Service
GOARCH = amd64

PLATFORM?=581
VERSION?=0.0.1

RELEASE_DIR=release/$(PLATFORM)


# Build flags
VPREFIX := gogit.oa.com/March/gopkg/build
GO_LDFLAGS := -X $(VPREFIX).Version=$(VERSION) -X $(VPREFIX).BuildUser=$(shell whoami)@$(shell hostname) -X $(VPREFIX).BuildDate=$(shell date +"%Y-%m-%dT%H:%M:%SZ")

# server local
.PHONY: default
default:
	cp cmd/server/app-example.toml release/server/app.toml
	go build -ldflags "$(GO_LDFLAGS)" -o release/server/$(BINARY) cmd/server/main.go

# client local
.PHONY: client
client:
	cp cmd/client/app-example.toml release/client/app.toml
	go build -ldflags "$(GO_LDFLAGS)" -o release/client/$(BINARY) cmd/client/main.go

.PHONY: dev
dev:
	mkdir -p $(RELEASE_DIR)
	cp cmd/server/conf/$(PLATFORM)/app-dev.toml $(RELEASE_DIR)/app.toml
	GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(GO_LDFLAGS)" -o $(RELEASE_DIR)/$(BINARY) cmd/server/main.go
	zip -jJ $(RELEASE_DIR)/$(BINARY)_v$(VERSION).zip $(RELEASE_DIR)/$(BINARY) $(RELEASE_DIR)/app.toml
