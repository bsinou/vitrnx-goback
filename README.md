# vitrnx-goback

[![GoDoc](https://godoc.org/github.com/bsinou/vitrnx-goback?status.svg)](https://godoc.org/github.com/bsinou/vitrnx-goback)
[![Build Status](https://travis-ci.org/bsinou/vitrnx-goback.svg?branch=master)](https://travis-ci.org/bsinou/vitrnx-goback)
 [![Go Report Card](https://goreportcard.com/badge/github.com/bsinou/vitrnx-goback)](https://goreportcard.com/report/github.com/bsinou/vitrnx-goback)

Simple Go backend for the vitrnx front.

It is also a sample project to experiment with Go language.

## Setup

### On CentOS

In order to run the vitrnx Go backend you need to:

```sh
# install and start mongo db
sudo yum install mongodb-server mongodb
sudo systemctl start mongod

# Insure that a data folder exists in one of the expected places:
mkdir -p ~/.config/vitrnx/<instanceID>/data
# Insure you have a correct firebase key file
cp firebase-apiCert.json ~/.config/vitrnx/<instanceID>/firebase-apiCert.json

# start the server
go run main.go start upala
```

### On Ubuntu

In order to run the vitrnx Go backend you need to:

```sh
# install and start mongo db
sudo apt install mongodb-server mongodb
sudo systemctl start mongod

# Insure that a data folder exists in one of the expected places:
mkdir -p ~/.config/vitrnx/<instanceID>/data
# Insure you have a correct firebase key file
cp firebase-apiCert.json ~/.config/vitrnx/<instanceID>/firebase-apiCert.json

# start the server
go run main.go start upala
```

Install on a CentOS server with Go.

```sh
cd go/src/github.com/bsinou/
git clone https://bsinou@github.com/bsinou/vitrnx-goback.git

go get firebase.google.com/go gopkg.in/mgo.v2/bson gopkg.in/mgo.v2 github.com/spf13/viper github.com/spf13/cobra github.com/mattn/go-sqlite3

cd vitrnx-goback/
make dev

sudo yum install mongodb-server mongodb
sudo systemctl start mongod
sudo systemctl enable mongod

sudo mkdir -p /etc/vitrnx/upala
sudo mkdir -p /var/lib/vitrnx/upala/data
sudo chown -R vitrnx.vitrnx /var/lib/vitrnx/upala/data/

# ON LOCAL WORKSTATION
scp  -P 2102 vitrnx.toml  bsinou@vps535638.ovh.net:/etc/vitrnx/upala

# Back on the server
cd go/src/github.com/bsinou/vitrnx-goback
go run main.go start upala

```