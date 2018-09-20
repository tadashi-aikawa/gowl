gowl
====

Go tools for GitHub and Bitbucket


Install
-------

TODO...


Usage
-----

```
$ gowl -h
Gowl.

Usage:
  gowl clone
  gowl edit
  gowl -h | --help
  gowl --version

Options:
  -h --help     Show this screen.
  --version     Show version.
```


Quick start
-----------

### Create `~/.gowlconfig`

`.gowlconfig` is a TOML file.

```toml
editor = "your editor command (ex: code)"

[github]
token = "your github token"
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
