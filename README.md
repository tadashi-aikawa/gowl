gowl
====

![](https://img.shields.io/github/release/tadashi-aikawa/gowl.svg)

An interactive git tool that works with GitHub and Bitbucket Server.

Support for both Windows and Linux.

![DEMO](https://raw.githubusercontent.com/tadashi-aikawa/gowl/master/demo.gif)

Install
-------

```
$ go get -u github.com/tadashi-aikawa/gowl
```

or

Download a binary from [release page](https://github.com/tadashi-aikawa/gowl/releases).


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

[editors]
default = "code"
vim = "vim"

[github]
token = "your github token"
# If `overrideUser = true`, Add userName and mailAddress to `.git/config` (`user.name` and `user.email`)
overrideUser = true
userName = "your github account name"
mailAddress = "your github email address"

[bitbucketserver]
baseurl = "http://your.bitbucket.server.url"
username = "yourname"
password = "yourpassword"
prefix = "your prefix in gopath (ex: mamansoft/bitbucket)"
useSSH = true
```

#### A minimum example

```toml
browser = "chrome"

[editors]
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

1. `gowl get`
2. `gowl edit`


Configuration
-------------

Gowl uses toml format as a configuration file.  
Please check `config.go`.

TODO: Definition table


Root directory
--------------

The root directory is determined by the following priority.

1. `root` in `.gowlconfig`
2. `<GOPATH>/src`
3. `<HOME>/.gowlroot`


Other
-----

If you use fzf(or peco), the following setting may make you happy!

bash
```
alias cdg="cd $(gowl list | fzf)"
```

fish
```
alias cdg "cd (gowl list | fzf)"
```

![DEMO2](https://raw.githubusercontent.com/tadashi-aikawa/gowl/master/demo2.gif)


For developer
-------------

### Prerequirements

* [dep](https://golang.github.io/dep/)


### Create environment

```
$ dep ensure
```

### Release

#### Requirements

* make
* bash
* dep
* ghr

#### Packaging and deploy

Confirm that your branch name equals release version, then...

```
$ make release
```
