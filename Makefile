# This how we want to name the binary output
BINARY=genframe

REPO=github.com/0xgenoskwa/gfsandbox

# These are the values we want to pass for VERSION and BUILD
VERSION:=1.0.0
TAG:=${shell git rev-parse --short HEAD}
BUILD:=${shell date +%FT%T%z}

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X ${REPO}/config.Version=${VERSION}-${TAG} -X ${REPO}/config.Build=${BUILD}"

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY} main.go

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
