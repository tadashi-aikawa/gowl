gowl
====

Go tools for GitHub and Bitbucket.

Support for both Windows and Linux.

![DEMO](https://raw.githubusercontent.com/tadashi-aikawa/gowl/master/demo.gif)

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
root = "Root directory of repositories for gowl"
# ex. C:\\users\\tadashi-aikawa\\.gowl
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

#### A minimum example

```toml
browser = "chrome"

[tools]
default = "code"

[github]
token = "your github token"
```

This file means...

* Use Google Chrome as browser
* Use VSCode as editor
* Use GitHub only


### Run

For example..

1. `gowl clone`
2. `gowl edit`


Root directory
--------------

The root directory is determined by the following priority.

1. `root` in `.gowlconfig`
2. `<GOPATH>/src`
3. `<HOME>/.gowlroot`


For developer
-------------

### Prerequirements

* [dep](https://golang.github.io/dep/)


### Create environment

```
$ dep ensure
```
