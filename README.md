# fsquota

[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/anexia-it/fsquota/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/anexia-it/fsquota?status.svg)](https://godoc.org/github.com/anexia-it/fsquota)
[![Build Status](https://travis-ci.org/anexia-it/fsquota.svg?branch=master)](https://travis-ci.org/anexia-it/fsquota)
[![codecov](https://codecov.io/gh/anexia-it/fsquota/branch/master/graph/badge.svg)](https://codecov.io/gh/anexia-it/fsquota)
[![Go Report Card](https://goreportcard.com/badge/github.com/anexia-it/fsquota)](https://goreportcard.com/report/github.com/anexia-it/fsquota)

fsquota is a native Go library for interacting with (Linux) filesystem quotas. This library does **not** make use of cgo
or invoke external commands, but rather interacts directly with the kernel interface by use of syscalls. This library is
maintained by the [Anexia](https://www.anexia-it.com/) R&D team.

## Portability

fsquota has been developed with Linux in mind and as such only supports Linux for now. Support for other platforms may
be added in the future.

## fsqm

This repository also ships *fsqm*, a simple command line interface to filesystem quotas. *fsqm* provides the ability to
retrieve user and group quota reports and management of user and group quotas.

*fsqm* can be obtained from [the releases page](https://github.com/anexia-it/fsquota/releases).

## fsqm(v0.1.3)

```shell
git clone https://github.com/anexia-it/fsquota.git && \
cd fsquota  && \
go mod init github.com/anexia-it/fsquota && \
go mod tidy && \
go mod vendor && \
cd cmd/fsqm && \
go build && \
go install
```

## How to use it

[Initialize System](https://anexia.com/blog/en/filesystem-quota-management-in-go/) [a sample file](https://github.com/anexia-it/wad2018-quotactl/blob/master/quotactl.go)

We can start with simple setup first

```shell
sudo su
truncate -s 1G /tmp/test.ext4 && \
/sbin/mkfs.ext4 /tmp/test.ext4 && \
mkdir -p /mnt/quota_test && \
mount -o usrjquota=aquota.user,grpjquota=aquota.group,jqfmt=vfsv1 /tmp/test.ext4 /mnt/quota_test && \
quotacheck -vucm /mnt/quota_test && \
quotacheck -vgcm /mnt/quota_test && \
quotaon -v /mnt/quota_test && \
lsblk
```

![lsblk should show all mouted items](https://i.imgur.com/aCQjKwM.png)

Thus, the file path will be `/dev/loop27`[/dev/loopN] from the initialize article.

We can run [a sample file](https://github.com/anexia-it/wad2018-quotactl/blob/master/quotactl.go)
get the sample file and run:

```shell
#            special:path userid
go run *.go "/dev/loop27" "0"
```

It will yeild the details

```notes
Space (1K Blocks):
  - hard limit: 0
  - soft limit: 0
  - usage     : 20480
Inodes:
  - hard limit: 0
  - soft limit: 0
  - usage     : 2
```

Now we can add limit using fsquota:

Docker Compile `make snapshot`

```shell
#                           mounted    username         softL,hardL
./bin/fsqm-amd64 user set "/dev/loop27" "a"     --files "1mb,2mb"
```

```notes
bytes:
  - soft: 0 B
  - hard: 0 B
  - used: 12 MiB
files:
  - soft: 1.0 M
  - hard: 2.0 M
  - used: 3
```

Which will create following files into `/mnt/quota_test` folder as following

```shell
sudo chmod -R 777 /mnt/quota_test # will give error ignore.
cd /mnt/quota_test && \
ls -lah
```

```shell
➜  /mnt cd /mnt/quota_test && \
ls -lah
total 13M
drwxrwxrwx 3 root root 4.0K Mar 16 22:53 .
drwxr-xr-x 4 root root 4.0K Mar 16 22:33 ..
-rw------- 1 root root 7.0K Mar 16 22:34 aquota.group
-rw------- 1 root root 7.0K Mar 16 22:34 aquota.user
drwxrwxrwx 2 root root  16K Mar 16 22:33 lost+found
```

Now create random files will not stop the user, not usre why developer may correct here:

```
dd if=/dev/urandom of=somefile.bin bs=1M count=1 && \
dd if=/dev/urandom of=somefile_x1.bin bs=1M count=1 && \
dd if=/dev/urandom of=somefile_x2.bin bs=1M count=1
```

## Issue tracker

Issues in fsquota are tracked using the corresponding GitHub
project's [issue tracker](https://github.com/anexia-it/fsquota/issues).

## Status

The current release is **v0.1.4**.

Changes to fsquota are subject to [semantic versioning](http://semver.org/).

## License

fsquota is licensed under the terms of the [MIT license](https://github.com/anexia-it/fsquota/blob/master/LICENSE).
