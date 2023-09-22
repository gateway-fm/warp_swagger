APP:=warp_swagger
APP_ENTRY_POINT:=warp.go
BUILD_OUT_DIR:=./
GOPRIVATE:=github.com

# Set GOOS and GOARCH to the current system values using the go env command
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# set git related vars for versioning
TAG 		:= $(shell git describe --abbrev=0 --tags)
COMMIT		:= $(shell git rev-parse HEAD)
BRANCH		?= $(shell git rev-parse --abbrev-ref HEAD)
REMOTE		:= $(shell git config --get remote.origin.url)
BUILD_DATE	:= $(shell date +'%Y-%m-%dT%H:%M:%SZ%Z')

# Set RELEASE to either the current TAG or COMMIT
RELEASE :=
ifeq ($(TAG),)
	RELEASE := $(COMMIT)
else
	RELEASE := $(TAG)
endif

# append versioner vars to ldflags
LDFLAGS += -X $(GITVER_PKG).ServiceName=$(APP)
LDFLAGS += -X $(GITVER_PKG).CommitTag=$(TAG)
LDFLAGS += -X $(GITVER_PKG).CommitSHA=$(COMMIT)
LDFLAGS += -X $(GITVER_PKG).CommitBranch=$(BRANCH)
LDFLAGS += -X $(GITVER_PKG).OriginURL=$(REMOTE)
LDFLAGS += -X $(GITVER_PKG).BuildDate=$(BUILD_DATE)

# The all target runs the tidy, build, and test targets
all: generate-proto tidy build test

# The tidy target runs go mod tidy
tidy:
	go mod tidy

# The update target runs go get -u
update:
	go get -u ./...


# the feature bellow allows to execute command
# without adding ARGS=... input
ifeq (generate,$(firstword $(MAKECMDGOALS)))
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS):;@:)
endif

doom:
	@MallocNanoZone=0 go run -race $(APP_ENTRY_POINT) dummy

summon:
	@MallocNanoZone=0 go run -race $(APP_ENTRY_POINT) summon
generate:
	@MallocNanoZone=0 go run -race $(APP_ENTRY_POINT) generate $(ARGS)
generate-proto:
	cd protocols/observer && protoc --go_out=paths=source_relative:. --go_opt=paths=source_relative  --go-grpc_out=paths=source_relative:. --go-grpc_opt=paths=source_relative  *.proto

update_subtree:
	git subtree split --rejoin --prefix protocols && git subtree pull --prefix=protocols git@github.com:gateway-fm/protocols.git main --squash

# no thi ng