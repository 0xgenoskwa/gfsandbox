# This how we want to name the binary output
BINARY=genframe

REPO=github.com/0xgenoskwa/gfsandbox

# These are the values we want to pass for VERSION and BUILD
GENFRAME_VERSION:=1.0.0
AUTOUPDATER_VERSION:=1.0.0
AUTORECOVER_VERSION:=1.0.0
TAG:=${shell git rev-parse --short HEAD}
BUILD:=${shell date +%FT%T%z}

# Setup the -ldflags option for go build here, interpolate the variable values
GENFRAME_LDFLAGS=-ldflags "-X ${REPO}/config.Version=${GENFRAME_VERSION}-${TAG} -X ${REPO}/config.Build=${BUILD}"
AUTOUPDAER_LDFLAGS=-ldflags "-X ${REPO}/config.Version=${GENFRAME_VERSION}-${TAG} -X ${REPO}/config.Build=${BUILD}"
AUTORECOVER_LDFLAGS=-ldflags "-X ${REPO}/config.Version=${GENFRAME_VERSION}-${TAG} -X ${REPO}/config.Build=${BUILD}"

generate:
	GOFLAGS=-mod=mod go generate ./...

# Builds the project
build-genframe:
	GOOS=linux GOARCH=amd64 go build ${GENFRAME_LDFLAGS} -o ${BINARY} cmd/genframe/main.go

build-autoupdater:
	GOOS=linux GOARCH=amd64 go build ${AUTOUPDAER_LDFLAGS} -o ${BINARY} cmd/genframe/main.go

build-autorecover:
	GOOS=linux GOARCH=amd64 go build ${AUTORECOVER_LDFLAGS} -o ${BINARY} cmd/genframe/main.go

test:
	go test ./...  -count=1 -cover -race

test-integration:
	go test ./... --tags=integration -count=1 -race

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
