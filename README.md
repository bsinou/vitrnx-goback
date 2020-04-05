# vitrnx-goback

[![GoDoc](https://godoc.org/github.com/bsinou/vitrnx-goback?status.svg)](https://godoc.org/github.com/bsinou/vitrnx-goback) [![Build Status](https://travis-ci.org/bsinou/vitrnx-goback.svg?branch=master)](https://travis-ci.org/bsinou/vitrnx-goback) [![Go Report Card](https://goreportcard.com/badge/github.com/bsinou/vitrnx-goback)](https://goreportcard.com/report/github.com/bsinou/vitrnx-goback)

Simple backend for the VitrnX front and also play and experiment with the **Go** language.

**IMPORTANT**
In the various scripts below, we expect you have define following variable, *on **both** the server and your local workstation* when relevant:

```sh
adminUser=sysadmin
myHost=XXX
myPort=XXX
appID=XXX
timestamp=$(date -u +%y%m%d)

# It helps to put this in a sh file and do:
source $GITSRCPATH/github.com/bsinou/vitrnx-goback/setup-variables.sh
```

## Overview

### In production

```sh
# Start the app
sudo systemctl restart $appID
# Data
cd /var/lib/vitrnx/$appID/data/
# Conf
vi /etc/vitrnx/$appID/vitrnx.toml
```

## Setup & Deploy

### Deploy update

```sh

## FROM SERVER

# On server stop instance
ssh $adminUser@$myHost -p $myPort
systemctl stop $appID

# if you want to push local DB to remote server, do first:
mkdir -p /var/lib/vitrnx/$appID/backup/$timestamp
mv /var/lib/vitrnx/$appID/data/gorm-sqlite.db /var/lib/vitrnx/$appID/backup/$timestamp

## FROM LOCAL WORKSTATION

# Update db content
cd ~/.config/vitrnx/$appID/data
scp -P $myPort gorm-sqlite.db $adminUser@$myHost:/var/lib/vitrnx/$appID/data/

# Update the binary
cd $GITSRCPATH/github.com/bsinou/vitrnx-goback
make prod
scp -P $myPort vitrnx-goback $adminUser@$myHost:/usr/local/bin

# Then on the server
systemctl restart $appID; journalctl -f -u $appID --since "2 minutes ago";
```

### Retrieve backups from the prod

```sh
ssh $adminUser@$myHost -p $myPort
mkdir -p /var/lib/vitrnx/$appID/backup/$timestamp
cd !$
#systemctl stop $appID
cp /var/lib/vitrnx/$appID/data/gorm-sqlite.db .
cp /etc/vitrnx/$appID/vitrnx.toml .
exit

# back on locahost
cd ~/.config/vitrnx/$appID/data
scp -P $myPort $adminUser@$myHost:/var/lib/vitrnx/$appID/backup/$timestamp/*.db .

```

### On CentOS

In order to run the vitrnx Go backend you need to:

```sh

# Insure that a data folder exists in one of the expected places:
mkdir -p ~/.config/vitrnx/$appID/data
# Insure configuration file is correct
# a sample file can be found in vitrnx-goback/conf sub folder
vi ~/.config/vitrnx/$appID/vitrnx.toml


# start the server
go run main.go start $appID
```

### On Ubuntu

In order to run the vitrnx Go backend you need to:

```sh
# Insure that a data folder exists in one of the expected places:
mkdir -p ~/.config/vitrnx/$appID/data
# OR
sudo mkdir -p /var/lib/vitrnx/$appID/data
sudo mkdir -p /etc/vitrnx/$appID
sudo chown -R <current user> /etc/vitrnx/$appID
sudo chown -R <current user> /var/lib/vitrnx/$appID

# start the server
go run main.go start $appID
```

## Setup workstation

Install on a CentOS server with Go.

```sh
cd $GITSRCPATH/github.com/bsinou/
git clone https://bsinou@github.com/bsinou/vitrnx-goback.git

cd vitrnx-goback/
make dev

sudo mkdir -p /etc/vitrnx/$appID
sudo mkdir -p /var/lib/vitrnx/$appID/data
sudo chown -R vitrnx.vitrnx /var/lib/vitrnx/$appID/data/

# ON LOCAL WORKSTATION
scp  -P $myPort vitrnx.toml  $adminUser@$myHost:/etc/vitrnx/$appID

# Back on the server... 
#Enhance, we don't want to install Go on the server
cd go/src/github.com/bsinou/vitrnx-goback
go run main.go start $appID
```

## Configure Cells as Oauth provider

Simply add this to your pydio.json file:

```json
{
    "client_id": "vitrnx",
    "client_name": "A simple blog manager",
    "grant_types": [
        "authorization_code",
        "refresh_token"
        ],
    "redirect_uris": [
        "http://localhost:3000/callback"
    ],
    "response_types": [
        "code",
        "token",
        "id_token"
    ],
    "scope": "openid email profile pydio offline"
}
```

## Developer Tips

To check a new http entry point.

```sh
curl http://localhost:8888/auth/ab-check
```

