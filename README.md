gowl
====

Go tools for GitHub and Bitbucket


Install
-------

TODO...


Usage
-----

TODO...


Quick start
-----------

### Create `.gowlconfig`

`.gowlconfig` is a TOML file.

```toml
[github]
token = "your github token"
```

### Run

```
$ gowl lodash
★556     29086069       JavaScript     lodash-archive/lodash-fp                                    2018-09-17 07:22:20 +0000 UTC
★415     36026026       TypeScript     steelsojka/lodash-decorators                                2018-09-18 09:18:03 +0000 UTC
★1       102501643      JavaScript     DevMountain/javascript-5-lodash                             2018-07-02 17:18:08 +0000 UTC
★112     56010441       JavaScript     jfmengels/eslint-plugin-lodash-fp                           2018-09-16 04:12:32 +0000 UTC
★64      38583433       JavaScript     mike-north/ember-lodash                                     2018-06-27 14:01:36 +0000 UTC
★8688    18351848       JavaScript     typicode/lowdb                                              2018-09-18 13:37:08 +0000 UTC
★104     12422224       JavaScript     lodash-archive/lodash-cli                                   2018-09-13 13:08:58 +0000 UTC
★211     17031775       JavaScript     marklagendijk/lodash-deep                                   2018-09-17 18:38:07 +0000 UTC
★104     15248146       HTML           node4good/lodash-contrib                                    2018-04-04 23:14:33 +0000 UTC
★319     29259882       CSS            davidkpiano/sassdash                                        2018-09-16 19:44:21 +0000 UTC
★116     100284075                     yeyuqiudeng/pocket-lodash                                   2018-09-04 08:31:59 +0000 UTC
★114     12422205       JavaScript     lodash-archive/lodash-node                                  2018-02-02 13:01:42 +0000 UTC
★2       64949810       HTML           learn-co-students/javascript-lodash-templates-v-000         2018-03-28 13:31:20 +0000 UTC
★13      53098262       TypeScript     types/npm-lodash                                            2018-09-16 04:06:16 +0000 UTC
★176     22221844       JavaScript     syzer/JS-Spark                                              2018-09-07 16:13:29 +0000 UTC
★561     36891832                      underdash/underdash                                         2018-09-16 04:07:19 +0000 UTC
★216     114683570      PHP            lodash-php/lodash-php                                       2018-09-14 06:32:57 +0000 UTC
★90      23645308       JavaScript     mtraynham/lodash-joins                                      2018-09-08 08:44:02 +0000 UTC
★267     37610237       JavaScript     cvgellhorn/webpack-boilerplate                              2018-09-15 10:26:13 +0000 UTC
```


For developer
-------------

### Prerequirements

* [dep](https://golang.github.io/dep/)


### Create environment

```
$ dep ensure
```
