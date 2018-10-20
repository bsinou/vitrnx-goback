# Resouces for Git

_Below tip and tricks are mainly intented to be use with Github but might also be relevant for raw git repos or services provided by others_.

## Configuration

### Commit using a ssh key

Thanks to [Developius](https://gist.github.com/developius/c81f021eb5c5916013dc)

Create a repo. Make sure there is at least one file in it (even just the README) Generate ssh key:

```sh
ssh-keygen -t rsa -C "your_email@example.com" -b 4096
```

Copy the contents of the file ~/.ssh/id_rsa.pub to your SSH keys in your GitHub account settings.

```sh
# Test SSH key:
ssh -T git@github.com

#clone the repo:
git clone git://github.com/username/your-repository
```

Now cd to your git clone folder and do:

```sh
git remote set-url origin git@github.com:username/your-repository.git

# Now try editing a file (try the README) and then do:
git add -A
git commit -am "my update msg"
git push
```

## Useful commands

### Cache your credentials

For linux:

```sh
# default: 15mn
git config --global credential.helper cache

# longer: set the cache to timeout after 1 hour (setting is in seconds)
git config --global credential.helper 'cache --timeout=3600'
```

For other distros please refer to [this blog post](https://help.github.com/articles/caching-your-github-password-in-git/)

## Various strategies to keep a fork up to date with the upstream server

- [A simple and concise tutorial](https://robots.thoughtbot.com/keeping-a-github-fork-updated)
- [playing with github api](https://medium.com/@durgaprasadbudhwani/playing-with-github-api-with-go-github-golang-library-83e28b2ff093)
- [go-github lib road map](https://docs.google.com/spreadsheets/d/1wlHnoJSAN01nXUUF1JobsaKO4mJjAwZNWK8L_4ZSlls/edit#gid=0)
