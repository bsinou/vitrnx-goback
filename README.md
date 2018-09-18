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