GOBUILD=go build
ENV=env GOOS=linux

TIMESTAMP:=$(shell date -u +%Y-%m-%dT%H:%M:%S)
GITREV:=$(shell git rev-parse HEAD)
VERSION?=0.8.3
XGO_TARGETS?="linux/amd64"

all: clean build

build: prod

prod:
	go build \
		-ldflags "-X github.com/bsinou/vitrnx-goback/conf.Env=production \
		-X github.com/bsinou/vitrnx-goback/conf.VitrnxVersion=${VERSION} \
		-X github.com/bsinou/vitrnx-goback/conf.BuildTimestamp=${TIMESTAMP} \
		-X github.com/bsinou/vitrnx-goback/conf.BuildRevision=${GITREV}" \
		-o vitrnx-goback main.go

dev:
	go build \
		-ldflags "-X github.com/bsinou/vitrnx-goback/conf.VitrnxVersion=${VERSION} \
		-X github.com/bsinou/vitrnx-goback/conf.BuildTimestamp=${TIMESTAMP} \
		-X github.com/bsinou/vitrnx-goback/conf.BuildRevision=${GITREV}" \
		-o vitrnx-goback main.go

cleanall: stop clean rm

clean:
	rm -f vitrnx-goback
