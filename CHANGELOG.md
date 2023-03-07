# CHANGELOG

## 0.3.3 - (2023-03-04)

### Fix

- 生成模型描述信息;
- 修改`List()`实现;
- 修改`MetaField.SchemaName`实现;

### Refactor

- `swag fmt`

## 0.3.2 - (2023-03-04)

### Feat

- 支持嵌套结构体文档的自动生成;

### TODO

- [x] 模型描述信息获取错误;

## 0.3.0 - (2023-03-02)

### BREAKING

- 重写`OpenApi`文档自动生成方案，引入`OpenApi`对象;
- 重写基本模型`BaseModel`，引入`godantic`包及系列接口;
- 修改`HandlerFunc`函数签名;
- 修改`CronJob`签名及实现;
- 删除原表示运行模式的包`mode`;
- 删除`zaplog`，`FlaskGo`的默认日志句柄由标准库`log`实现;
- 删除`FlaskGo`的控制台日志输出，仅保留`Logger()`接口;

## 0.2.9 - (2023-01-01)

### Feat

- 新增`redoc`文档页面;

### Fix

- 修复`Router.IncludeRouter`覆盖同名路由的错误;

## 0.2.8 - (2022-11-01)

### Fix

- 修复当struct的json字段标签存在多个选项时，swagger未正确提取字段名的错误;

## 0.2.7 - (2022-10-31)

### Refactor

- update `functools` to `v0.1.6`;
- 代码格式化;
- `struct`内存对齐;

## 0.2.6 - (2022-10-22)

### Fix

- 自定义请求错误返回`400`状态码;

## 0.2.5 - (2022-10-22)

### Feat

- 路径参数校验可用；
- 查询参数校验可用；
- 升级`functools`到`v0.1.5`;

### Refactor

- `QModel`新增`InPath`参数，适用于查询参数和路径参数;
- 分类存储`FlaskGo.Route`;

### Supported Features

- [x] swagger 文档自动生成；
- 请求体自动校验：
    - [x] 路径参数自动校验；
    - [x] 查询参数自动校验；
    - [ ] 请求体自动校验, 需通过`Context.ShouldBindJSON()`主动校验；
- [ ] 响应体自动校验；
- [ ] 平滑关机；

## 0.2.4 - (2022-10-21)

### Feat

- 新增方法`ReplaceErrorHandler`, 以支持替换默认的fiber错误处理方法；

## 0.2.3 - (2022-10-21)

### Refactor

- `Context`现在关联了各种`Response`方法，无需再通过`flaskgo`调用；

## 0.2.2 - (2022-10-20)

### Fix

- 支持识别可选路径参数和必选路径参数；
- `Route`新增路径参数,支持路径参数校验;

## 0.2.1 - (2022-10-19)

### Fix

- 修复当接口请求体，响应体为nil时，无法生成文档的错误；
- 排除空匿名结构体, 避免`flaskgo.BaseModel`出现在文档中；

### Refactor

- 修改`RouteModel.Struct`为`BaseModelIface`类型；

### Supported Features

- [x] swagger 文档自动生成；
- [ ] 请求体自动校验；
- [ ] 响应体自动校验；
- [ ] 平滑关机；

## 0.2.0 - (2022-10-19)

### BREAKING

- 修改`APIRouter`接口参数类型，新增`BaseModelIface`接口，请求体和相应体均需实现此接口;

## 0.1.6 - (2022-10-19)

### Feat

- swagger 自动生成可用；

### Supported Features

- [x] swagger 文档自动生成；
- [ ] 请求体自动校验；
- [ ] 响应体自动校验；
- [ ] 平滑关机；

## 0.1.5 - (2022-10-18)

### Feat

- 新增方法：`Route.AddDependencies`;

### Refactor

- 修改事件注册方式：`FlaskGo.OnEvent()`;
- 默认`禁用多进程`;

### Rename

- `FlaskGo.Context` -> `FlaskGo.Service`;

## 0.1.4 - (2022-10-17)

### Feat

- 新增`AdvancedResponse`返回体;
- 修改`AnyResponse`的返回实现;

## 0.1.3 - (2022-10-17)

### Feat

- 修改文件组织结构；
- 修正swagger文档错误；

## 0.1.2 - (2022-10-15)

### Feat

- swagger文档页面可见（未测试）；

## 0.1.1 - (2022-10-15)

### Feat

- 修改模型自动校验（未测试）；

## 0.1.0 - (2022-10-14)

### Feat

- 提供对`struct`的校验方法和错误信息适配；
- 修改`Context`;
- 升级`functools`到`v0.1.4`;

### Rename

- `Context` -> `Service`;
- new `Context`

### Bug

- 多进程模式下程序无法允许；

## 0.0.5 - (2022-10-13)

### Feat

### Refactor

- 修改`fiber`请求日志格式；

## 0.0.4 - (2022-10-11)

### Feat

- 移除`logger`,引入`functools/zaplog`;

### Refactor

- 格式化代码;

## 0.0.3 - (2022-10-11)

### Feat

- `Fiber`现在正常工作;
- `FlaskGoContext` 新增表单校验方法（未完成）；
- 为`Fiber`注册`recover`中间件，`ErrorsHandler`中间件、日志中间件；
- 重写swagger生成方法（未完成）；
- 修改实例项目；
- API基本固定；

### Refactor

- 默认开启请求体自动校验；

### TODO

- swagger 现在不能正常工作;
- 请求体自动校验现在不能正常工作;

## 0.0.2 - (2022-10-09)

### Feat

- 移除基础组件并引入软件包：`functools`;

### Rename

- `FlaskGoContext` -> `Context`;
- `FlaskGoResponse` -> `Response`;
- `FlaskGoRouter` -> `Router`;
- `FlaskGoRoute` -> `Route`;

## 0.0.1 - (2022-10-08)

### INIT

- 切换底层`gin`为`fiber`;