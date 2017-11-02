# rotatee

[![Build Status](https://travis-ci.org/saka1/rotatee.svg?branch=master)](https://travis-ci.org/saka1/rotatee)

rotatee is a simple program that works like tee with some extension:

- rotatee copies input to output file(s) with rotation.
- Rotation is period-based

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
$ echo "abcdefg" | ./rotatee --size 3B /tmp/size%i
abcdefg
$ cat /tmp/size1 /tmp/size2 /tmp/size3
abcdefg
```


## TODO
- Add more test cases
- Size-based rolling
- Max-history(limit number of files to keep)
- Rolling with archive
- Improve some inconsistent behavior