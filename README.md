# fsquota

[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/anexia-it/fsquota/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/anexia-it/fsquota?status.svg)](https://godoc.org/github.com/anexia-it/fsquota)
[![Build Status](https://travis-ci.org/anexia-it/fsquota.svg?branch=master)](https://travis-ci.org/anexia-it/fsquota)
[![codecov](https://codecov.io/gh/anexia-it/fsquota/branch/master/graph/badge.svg)](https://codecov.io/gh/anexia-it/fsquota)
[![Go Report Card](https://goreportcard.com/badge/github.com/anexia-it/fsquota)](https://goreportcard.com/report/github.com/anexia-it/fsquota)


fsquota is a native Go library for interacting with (Linux) filesystem quotas.
This library does **not** make use of cgo or invoke external commands, but rather interacts directly with the kernel interface by use of syscalls.
This library is maintained by the [Anexia](https://www.anexia-it.com/) R&D team.

## Portability

fsquota has been developed with Linux in mind and as such only supports Linux for now.
Support for other platforms may be added in the future.

## fsqm

This repository also ships *fsqm*, a simple command line interface to filesystem quotas. *fsqm* provides the ability to retrieve user and group quota reports and management of user and group quotas.

*fsqm* can be obtained from [the releases page](https://github.com/anexia-it/fsquota/releases).

## Issue tracker

Issues in fsquota are tracked using the corresponding GitHub project's [issue tracker](https://github.com/anexia-it/fsquota/issues).

## Status

The current release is **v0.1.3**.


Changes to fsquota are subject to [semantic versioning](http://semver.org/).

## License

fsquota is licensed under the terms of the [MIT license](https://github.com/anexia-it/fsquota/blob/master/LICENSE).
