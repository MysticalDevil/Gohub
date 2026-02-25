# Gohub
基于 Gin 框架的 Go 语言论坛 API 项目

# 使用到的开源库

- [`gin`](https://github.com/gin-gonic/gin) 路由、路由组、中间件

- `log/slog` 标准库日志

- [`gorm`](https://github.com/go-gorm/gorm) ORM 数据操作

- [`cobra`](https://github.com/spf13/cobra) 命令行结构

- [`viper`](https://github.com/spf13/viper) 配置信息

- [`cast`](https://github.com/spf13/cast) 类型转换

- [`redis`](https://github.com/redis/go-redis) Redis 操作

- [`jwt`](https://github.com/golang-jwt/jwt) JWT 操作

- [`base64Captcha`](https://github.com/mojocn/base64Captcha) 图片验证码

- [`validator`](https://github.com/go-playground/validator) 请求验证器

- [`limiter`](https://github.com/ulule/limiter) 限流器

- [`email`](https://github.com/jordan-wright/email) SMTP 邮件发送

- [`ansi`](https://github.com/mgutz/ansi) 终端高亮输出

- [`strcase`](https://github.com/iancoleman/strcase) 字符串大小写操作

- [`pluralize`](https://github.com/gertd/go-pluralize) 英文字符单数复数处理

- [`faker`](https://github.com/go-faker/faker) 假数据填充

- [`imaging`](https://github.com/disintegration/imaging) 图片裁切

# 自定义包

- app 应用对象
- auth 用户授权
- cache 缓存
- captcha 图片验证码
- config 配置信息
- console 终端
- database 数据库操作
- file 文件处理
- hash 哈希
- helpers 辅助方法
- jwt JWT 认证
- limiter API 限流
- logger 日志处理
- mail 邮件发送
- migrate 数据库迁移
- paginator 分页器
- redis Redis 数据库操作
- response 响应处理
- seed 假数据填充
- sms 发送短信
- str 字符串处理
- verifycode 数字验证码

# 支持的指令

```shell
$ go run main.go -h
Default will run "serve" command, you can use "-h" flag to see all subcommands

Usage:
  Gohub [command]

Available Commands:
  cache       Cache management
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  key         Generate App Key, will print the generated key
  make        Generate file nad code
  migrate     Run database migration
  play        Likes the Go Playground, but running at our application context
  seed        Insert fake data to the database
  serve       Start web server

Flags:
  -e, --env string   load .env file, example: --env=testing will use .env.testing file
  -h, --help         help for Gohub

Use "Gohub [command] --help" for more information about a command.
```

## 配置提示
`APP_KEY` 必须是安全随机值。可通过 `go run main.go key` 生成并填入 `.env`。
`APP_ENV_PATH` 可指定自定义 env 文件路径（例如测试场景），优先级高于 `-e/--env` 与默认 `.env`。
测试中若设置 `CONSOLE_SILENT=1`，将静默控制台输出（仅在 `APP_ENV=testing` 时生效）。

## API 响应格式
统一响应结构：
```json
{
  "code": "OK",
  "msg": "OK",
  "data": {}
}
```

校验/处理失败：
```json
{
  "code": "ERR_VALIDATION",
  "msg": "Request verification failed, please see errors for details",
  "errors": {
    "field": ["message"]
  }
}
```

分页使用 `offset/limit`，返回：
```json
{
  "code": "OK",
  "msg": "OK",
  "data": {
    "items": [],
    "offset": 0,
    "limit": 20,
    "total": 200
  }
}
```

# TODO
Postman 文档书写
支持多种缓存中间件，目前只支持 Redis
使用多种 Web 框架重构，例如 Iris，Fiber 等
