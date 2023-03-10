# flaskgo

类Flask风格的基于`fiber`封装的参考`FastApi`实现的模板

## Usage:

### 查看在线文档

```bash
# 安装godoc
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060

# 或：pkgsite 推荐
go install golang.org/x/pkgsite/cmd/pkgsite@latest
cd flaskgo/
pkgsite -http=:6060 -list=false
# 浏览器打开：http://127.0.0.1:6060/github.com/Chendemo12/flaskgo
```

### `struct`内存对齐

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

fieldalignment -fix ./... 
```

### 导入`functools`

- Linux

```bash
go get -v github.com/Chendemo12/functools
```

- Windows

```powershell
# $Env:GOSUMDB = "off"
go get -v github.com/Chendemo12/flaskgo
```

### 作为第三方库导入

```bash
go get https://github.com/Chendemo12/flaskgo
```

## Examples:

### Guide

- [guide example](./test/example.go)

## TODO:

- [x] 平滑关机；
- swagger文档自动生成：
    - [x] 支持数组类型的返回值文档生成；
    - [x] 支持匿名结构体字段；
    - [x] 支持结构体嵌套数组的文档生成；
    - [x] 结构体嵌套结构体的文档生成；
- 请求体自动校验
    - [x] 请求体校验方法；
- 响应体：
    - [x] 支持响应字节流；

## 一些常用的API

- 全部`api`可见[`alias.go`](./alias.go)文件；
