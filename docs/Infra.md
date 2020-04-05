# Infrastructure resources

## Install on CentOS

### Caddy

On May 1st 2018, I followed [this post](https://www.hugeserver.com/kb/install-caddy-centos-7/) that did the trick on CentOS7.

**Note**: before starting caddy for the first time:

- insure your DNS is correctly set and known
- do not forget to configure firewall.

You should use the Let's Encrypt staging environment to test your setup and insure everything is set correctly before trying in the production environment.
If the Let's Encrypt registration process fail 5 times, you have to wait for an hour to try again.  

```sh
curl -s https://getcaddy.com | bash -s personal
# insure every thing is OK
which caddy
# Set correct permission (or do the first line with sudo)
sudo chown root:root /usr/local/bin/caddy
sudo chmod 755 /usr/local/bin/caddy
# enable binding with protected ports
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/caddy

# Manage group and permission
sudo mkdir -p /var/www/caddy
sudo groupadd caddy
sudo useradd -g caddy --home-dir /var/www/caddy --no-create-home --shell /usr/sbin/nologin --system caddy
sudo chown -R caddy.caddy /var/www/caddy

sudo chown -R root:caddy /etc/caddy
sudo mkdir /etc/ssl/caddy
chown -R caddy:root /etc/ssl/caddy
sudo chown -R caddy:root /etc/ssl/caddy
sudo chmod 770 /etc/ssl/caddy

# Create and populate a caddy file
sudo touch /etc/caddy/Caddyfile
sudo chown -R caddy.root /etc/ssl/caddy
sudo chown caddy.caddy /etc/caddy/Caddyfile
sudo chmod 444 /etc/caddy/Caddyfile
sudo chmod -R 555 /var/www/caddy

# Configure caddy as a service
sudo vi /etc/systemd/system/caddy.service
sudo chown root:root /etc/systemd/system/caddy.service
sudo chmod 644 /etc/systemd/system/caddy.service
sudo systemctl daemon-reload
sudo systemctl enable caddy

# configure firewall
sudo firewall-cmd --permanent --zone=public --add-service=http
sudo firewall-cmd --permanent --zone=public --add-service=https
sudo firewall-cmd --reload
```

## Deploy

### On a brand new server

Create tree structure

```sh
# Insure to have the certificate to communicate with firebase
sudo mkdir -p /var/lib/vitrnx/conf
sudo mv firebase-apiCert.json /var/lib/vitrnx/conf

# Add a log folder

```

## Legacy

Below are a few hints from the time when we first experimented with MongoDB as backend for the blog post and Google firebase as OAuth provider.

## Retrieve API Cert for firebase

to add Firebase to your app, please follow instruction found on [firebase documentation website](https://firebase.google.com/docs/admin/setup).

>To use the Firebase Admin SDKs, you'll need a Firebase project, a service account to communicate with the Firebase service, and >a configuration file with your service account's credentials.
>
> - If you don't already have a Firebase project, add one in the Firebase console.
> - Navigate to the Service Accounts tab in your project's settings page.
> - Click the Generate New Private Key button at the bottom of the Firebase Admin SDK section of the Service Accounts tab.

## Install MongoDB

Main problem with MongoDB is that it necessitates a *lot* of resources, even when idle.
It is then a waste of resources for a small blog that only serves a few requests (up until at least ~1000) a day.
We then decided to rather use a single SQL lite instance.

Anyway, you might find here the instructions [to install on Ubuntu 18.04 - Bionic Beaver](https://linuxconfig.org/how-to-install-latest-mongodb-on-ubuntu-18-04-bionic-beaver-linux).

On CentOS 7+, simply do:

```sh
sudo yum install mongodb-server mongodb
sudo systemctl start mongod
sudo systemctl enable mongod
```
