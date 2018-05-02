GOBUILD=go build
ENV=env GOOS=linux

all: clean build

build: main

main:
	go build -ldflags "-X github.com/bsinou/vitrnx-goback/conf.Env=production -X github.com/bsinou/vitrnx-goback/conf.VitrnxVersion=0.1.1 -X github.com/bsinou/vitrnx-goback/conf.BuildTimestamp=`date -u +%Y-%m-%dT%H:%M:%S` -X github.com/bsinou/vitrnx-goback/conf.BuildRevision=`git rev-parse HEAD`" -o vitrnx-goback main.go

dev:
	go build -ldflags "-X github.com/bsinou/vitrnx-goback/conf.VitrnxVersion=0.1.1 -X github.com/bsinou/vitrnx-goback/conf.BuildTimestamp=`date -u +%Y-%m-%dT%H:%M:%S` -X github.com/bsinou/vitrnx-goback/conf.BuildRevision=`git rev-parse HEAD`" -o vitrnx-goback main.go

cleanall: stop clean rm

clean:
	rm -f vitrnx-goback
