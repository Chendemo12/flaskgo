# flaskgo

类Flask风格的基于gin封装的参考FastApi实现的模板

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
# 浏览器打开：http://127.0.0.1:6060/gitlab.cowave.com/gogo/flaskgo
```

### `struct`内存对齐

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

fieldalignment -fix ./... 
```

### 导入`functools`

- Linux

```bash
echo "machine gitlab.cowave.com login chenguang.li@cowave.com password glpat-zVc6XisKQfy-TPxzCym8" > ~/.netrc

export GOSUMDB=off
go get -v gitlab.cowave.com/gogo/functools
```

- Windows

```powershell
echo "machine gitlab.cowave.com login chenguang.li@cowave.com password glpat-zVc6XisKQfy-TPxzCym8" > ~/.netrc

$Env:GOSUMDB = "off"
go get -v gitlab.cowave.com/gogo/flaskgo
```

### 作为第三方库导入

```bash
go get https://gitlab.cowave.com/gogo/flaskgo
```

## Examples:

### Guide

- [guide example](./test/example.go)
- [sample](https://gitlab.cowave.com/gogo/sample)

## TODO:

- [ ] 平滑关机；
- swagger文档自动生成：
    - [x] 支持数组类型的返回值文档生成；
    - [ ] 支持匿名结构体字段；
    - [x] 支持结构体嵌套数组的文档生成；
    - [x] 结构体嵌套结构体的文档生成；
- 请求体自动校验
    - [ ] 请求体校验方法；
- 响应体：
    - [x] 支持响应字节流；

## 一些常用的API

- 全部`api`可见[`alias.go`](./alias.go)文件；

#### `FlaskGo`相关方法：

- `NewFlaskGo`:
- `GetFlaskGo`:
- `Response`:
- `JSONResponse`:
- `ResourceNotFound`:
- `APIRouter`:
- `AddResponseHeader`:
- `DeleteResponseHeader`:

结构体相关方法：

- `StructToMap`:
- `StructReflect`:
- `GetStructName`:
- `GetStructFullName`:
- `GetStructFieldsValue`:
- `GetStructFieldsName`:
- `GetStructFieldsTags`:
- `GetStructFieldsType`:
- `StructToJson`:
- `StructToString`:
- `FormatStruct`:
- `MakeStruct`:
