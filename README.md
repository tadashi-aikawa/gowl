gowl
====

Go tools for GitHub and Bitbucket


Install
-------

```
$ go get -u github.com/tadashi-aikawa/gowl
```


Usage
-----

```
$ gowl --help
Gowl.

Usage:
  gowl get [--bitbucket-server] [-f]
  gowl edit <tool>
  gowl web
  gowl list
  gowl -h | --help
  gowl --version

Options:
  <tool>                      Use tool
  -B --bitbucket-server       Use Bitbucket Server
  -f --force                  Force remove and reclone if exists
  -h --help                   Show this screen.
  --version                   Show version.
```


Quick start
-----------

### Create `~/.gowlconfig`

`.gowlconfig` is a TOML file.

```toml
browser = "Your browser"
# ex. C:\\Program Files (x86)\\Google\\Chrome\\Application\\Chrome.exe

[tools]
code = "code"
vim = "vim"

[github]
token = "your github token"

[bitbucketserver]
baseurl = "http://your.bitbucket.server.url"
username = "yourname"
password = "yourpassword"
prefix = "your prefix in gopath (ex: mamansoft/bitbucket)"
```


### Run

For example..

1. `gowl clone`
2. `gowl edit`


For developer
-------------

### Prerequirements

* [dep](https://golang.github.io/dep/)


### Create environment

```
$ dep ensure
```
