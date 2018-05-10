GOBUILD=go build
ENV=env GOOS=linux

CurrentVersion=0.1.2

all: clean build

build: main

main:
	go build -ldflags "-X github.com/bsinou/vitrnx-goback/conf.Env=production -X github.com/bsinou/vitrnx-goback/conf.VitrnxVersion=$(CurrentVersion) -X github.com/bsinou/vitrnx-goback/conf.BuildTimestamp=`date -u +%Y-%m-%dT%H:%M:%S` -X github.com/bsinou/vitrnx-goback/conf.BuildRevision=`git rev-parse HEAD`" -o vitrnx-goback main.go

dev:
	go build -ldflags "-X github.com/bsinou/vitrnx-goback/conf.VitrnxVersion=$(CurrentVersion) -X github.com/bsinou/vitrnx-goback/conf.BuildTimestamp=`date -u +%Y-%m-%dT%H:%M:%S` -X github.com/bsinou/vitrnx-goback/conf.BuildRevision=`git rev-parse HEAD`" -o vitrnx-goback main.go

cleanall: stop clean rm

clean:
	rm -f vitrnx-goback
