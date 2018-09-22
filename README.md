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
  gowl get [-s | --shallow] [-f | --force] [-B | --bitbucket-server]
  gowl edit [-e=<editor> | --editor=<editor>]
  gowl web
  gowl list
  gowl -h | --help
  gowl --version

Options:
  -e --editor=<editor>        Use editor [default: default]
  -s --shallow                Use shallow clone
  -f --force                  Force remove and reclone if exists
  -B --bitbucket-server       Use Bitbucket Server
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
default = "code"
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
