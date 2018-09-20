gowl
====

Go tools for GitHub and Bitbucket


Install
-------

TODO...


Usage
-----

```
$ gowl --help
Gowl.

Usage:
  gowl repo <word>
  gowl repo <word> clone [<seq>]
  gowl edit <word> [<seq>]
  gowl -h | --help
  gowl --version

Options:
  <word>        Search word for repository.
  <seq>         Specify selections
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

Search repositories by word.

```
$ gowl repo jumeaux
1    *2       26078355       Python         tadashi-aikawa/jumeaux                                      2018-09-07 14:14:41 +0000 UTC
2    *0       25208664                      aragnom/twinsurf                                            2014-10-14 13:38:58 +0000 UTC
3    *0       102268648      Shell          tadashi-aikawa/jumeaux-toolbox                              2017-12-10 14:00:16 +0000 UTC
4    *0       80344324       TypeScript     tadashi-aikawa/miroir                                       2018-09-13 04:24:48 +0000 UTC
```

Then clone first repository, `tadashi-aikawa/jumeaux`.  
Only you have to do is adding `clone` to tail.

```
$ gowl repo jumeaux clone
Clone https://github.com/tadashi-aikawa/jumeaux.git to C:\Users\syoum\Go/src/github.com/tadashi-aikawa/jumeaux
```

`C:\Users\syoum\Go` is my `$GOPATH`.

You can also clone not only first repository but also others, by specify a number.  

```
$ gowl repo jumeaux clone 3
Clone https://github.com/tadashi-aikawa/jumeaux-toolbox.git to C:\Users\syoum\Go/src/github.com/tadashi-aikawa/jumeaux-toolbox
```


For developer
-------------

### Prerequirements

* [dep](https://golang.github.io/dep/)


### Create environment

```
$ dep ensure
```
