[![Build Status](https://github.com/axetroy/dvs/workflows/test/badge.svg)](https://github.com/axetroy/dvs/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/dvs)](https://goreportcard.com/report/github.com/axetroy/dvs)
![Latest Version](https://img.shields.io/github/v/release/axetroy/dvs.svg)
![License](https://img.shields.io/github/license/axetroy/dvs.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/dvs.svg)

## Docker-based Virtual System

`dvs` is a command-line tool for creating an isolated sandbox

Required: Docker

Features:

- [x] Cross-platform support
- [x] Creating an isolated Linux environment based on a sandbox
- [x] No residue. The container is deleted every time the process exits, don't worry it will fill up your disk space

### Usage

```bash
# linux repl
$ dvs

# run command in the linux sandbox
$ dvs run ls -lh

# Run the specified sandbox with Docker Image's name
$ dvs --image node:latest run node --version
```

### Installation

If you are using Linux/macOS. you can install it with the following command:

```shell
# install latest version
curl -fsSL https://raw.githubusercontent.com/axetroy/dvs/master/install.sh | bash
# or install specified version
curl -fsSL https://raw.githubusercontent.com/axetroy/dvs/master/install.sh | bash -s v0.1.0
```

Or

Download the executable file for your platform at [release page](https://github.com/axetroy/dvs/releases)

Then set the environment variable.

eg, the executable file is in the `~/bin` directory.

```bash
# ~/.bash_profile
export PATH="$PATH:~/bin"
```

finally, try it out.

```bash
dvs --help
```

### Upgrade

You can re-download the executable and overwrite the original file.

or run the following command to upgrade

```bash
$ dvs upgrade # upgrade to latest
$ dvs upgrade v0.2.0 # Update to specified version
```

### Uninstall

just delete `dvs` executable file

### Build from source code

Make sure you have `Golang@v1.13.1` installed.

```shell
$ git clone https://github.com/axetroy/dvs.git $GOPATH/src/github.com/axetroy/dvs
$ cd $GOPATH/src/github.com/axetroy/dvs
$ make build
```

### Test

```bash
$ make test
```

### License

The [MIT License](LICENSE)
