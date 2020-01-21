[English](README.md) | 中文简体

[![Build Status](https://github.com/axetroy/dvs/workflows/test/badge.svg)](https://github.com/axetroy/dvs/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/dvs)](https://goreportcard.com/report/github.com/axetroy/dvs)
![Latest Version](https://img.shields.io/github/v/release/axetroy/dvs.svg)
![License](https://img.shields.io/github/license/axetroy/dvs.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/dvs.svg)

## 基于 Docker 的虚拟环境系统

`dvs` 是用于创建隔离沙箱的命令行工具

使用前提: Docker

特性:

- [x] 跨平台支持
- [x] 创建隔离的 Linux 沙盒环境
- [x] 无残留。每当进程退出时，容器都会被删除，不用担心它会填满您的磁盘空间

### 使用方法

```bash
# 运行 Linux 的 repl
$ dvs

# 在 Linux 沙盒中运行命令
$ dvs run ls -lh

# 运行指定的 Docker 镜像
$ dvs --image node:latest run node --version
```

### 安装

在 [release page](https://github.com/axetroy/dvs/releases) 页面下载你平台相关的可执行文件

然后设置环境变量

例如, 可执行文件放在 `~/bin` 目录

```bash
# ~/.bash_profile
export PATH="$PATH:~/bin"
```

最后，试一下是否设置正确

```bash
dvs --help
```

### 升级

你可以重新下载可执行文件然后覆盖

或者输入以下命令进行升级到最新版

```bash
> dvs upgrade
```

### 从源码构建

```bash
> go get -v -u github.com/axetroy/dvs
> cd $GOPATH/src/github.com/axetroy/dvs
> make build
```

### 测试

```bash
make test
```

### 开源许可

The [MIT License](LICENSE)
