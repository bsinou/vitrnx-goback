# Infrastructure resources 

## Install on CentOS

### Install MongoDB

```sh
sudo yum install mongodb-server mongodb
sudo systemctl start mongod
sudo systemctl enable mongod
```

### Caddy

On May 1st 2018, I followed [this post](https://www.hugeserver.com/kb/install-caddy-centos-7/) that did the trick on CentOS7.

**Note**: before starting caddy for the first time:

- insure your DNS is correctly set and known
- do not forget to configure firewall.

Otherwise letsencrypt registration process will fail 5 times and you will have to wait for an hour to try again. You might also use the ltsencrypt staging env during this phase to insure everything is set correctly before trying on the prod environment.

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

