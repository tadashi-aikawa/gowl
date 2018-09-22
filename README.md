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
$ gowl -h
Gowl.

Usage:
  gowl get [--bitbucket-server]
  gowl edit <editor>
  gowl -h | --help
  gowl --version

Options:
  <editor>                    Use editor
  -B --bitbucket-server       Use Bitbucket Server
  -h --help                   Show this screen.
  --version                   Show version.
```


Quick start
-----------

### Create `~/.gowlconfig`

`.gowlconfig` is a TOML file.

```toml
[editors]
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

`gowl clone` or `gowl edit`


For developer
-------------

### Prerequirements

* [dep](https://golang.github.io/dep/)


### Create environment

```
$ dep ensure
```
