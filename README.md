# rotatee

[![Build Status](https://travis-ci.org/saka1/rotatee.svg?branch=master)](https://travis-ci.org/saka1/rotatee)
[![Go Report Card](https://goreportcard.com/badge/github.com/saka1/rotatee)](https://goreportcard.com/report/github.com/saka1/rotatee)
[![Coverage Status](https://coveralls.io/repos/github/saka1/rotatee/badge.svg?branch=master&v=2.0)](https://coveralls.io/github/saka1/rotatee?branch=master)

rotatee is a simple program that works like tee with some extension:

- rotatee copies input to output file(s) with rotation.
- Multiple rotation methods supported
  - size-based
  - period-based

rotatee is inspired by some popular programs such as 

- [rotatelogs](https://httpd.apache.org/docs/2.4/programs/rotatelogs.html)
- [cronolog](http://linux.die.net/man/1/cronolog)
- [logback](http://logback.qos.ch/)

## Install
TBD

## Options
See `rotatee --help`.

## Examples
Like tee:
```
$ echo 'hoge' | rotatee hoge.txt
```

Rolling period is guessed from the format:
```bash
$ for i in `seq 1 3`; do echo 'hoge'; sleep 1; done | rotatee /tmp/hoge-%S.log
hoge
hoge
hoge
$ ls /tmp/ | grep hoge
hoge-10.log
hoge-11.log
hoge-12.log
hoge-13.log
```

Size-based rolling:
```bash
$ echo "abcdefg" | ./rotatee --size 3B --history 10 /tmp/size%i.log
abcdefg
$ cat /tmp/size2.log /tmp/size1.log /tmp/size.log
abcdefg
```


## TODO
- Add more test cases
- Rolling with archive
