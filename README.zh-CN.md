# gin-generator/http

[English](./README.md) | 中文

一个用于 Gin 的轻量请求解析辅助库：先绑定请求体数据，再把 `path`/`query` 参数解析到同一个结构体中。

## 功能特性

- 先通过 `c.ShouldBind` 绑定 JSON/form 数据，通过 `json` 标签支持任意层嵌套的 struct、map、slice。
- 在顶层标量字段上支持自定义标签：
  - `path:"name"`：路由参数
  - `query:"name"`：查询参数
- `path`/`query` 支持的标量类型：
  - `string`
  - 有符号整数（`int`、`int8`、`int16`、`int32`、`int64`）
  - 无符号整数（`uint`、`uint8`、`uint16`、`uint32`、`uint64`）
  - `bool`
  - 浮点数（`float32`、`float64`）
- 使用 `validate` 标签自动校验，基于 [go-playground/validator](https://github.com/go-playground/validator)。

## 安装

```bash
go get github.com/gin-generator/http
```

## 使用示例

```go
package main

import (
    "net/http"

    _http "github.com/gin-generator/http"
    "github.com/gin-gonic/gin"
)

type Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type CreateOrderRequest struct {
    UserID  int               `path:"user_id"`
    Page    int               `query:"page" validate:"required,min=1"`
    Name    string            `json:"name" validate:"required"`
    Phone   string            `json:"phone" validate:"phone"`
    Tags    []string          `json:"tags"`
    Meta    map[string]string `json:"meta"`
    Address Address           `json:"address"`
    Items   []Address         `json:"items"`
}

func main() {
    r := gin.Default()

    r.POST("/users/:user_id/orders", func(c *gin.Context) {
        var req CreateOrderRequest
        if err := _http.Parse(c, &req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, req)
    })

    _ = r.Run(":8080")
}
```

## 工作原理

1. `Parse` 先调用 `c.ShouldBind(obj)`，绑定所有 `json`/`form` 标签字段，支持任意深度的嵌套 struct、map、slice。
2. 使用反射遍历顶层结构体字段。
3. 对每个标量字段读取 `path`/`query` 标签，将字符串参数转换为对应字段类型。
4. 非标量字段（`struct`、`map`、`slice` 等）直接跳过——`ShouldBind` 已完成填充。
5. 最后使用全局 validator 实例对整个结构体进行 `validate` 标签校验。

## 标签行为说明

| 标签 | 数据来源 | 适用范围 |
|---|---|---|
| `json:"name"` | 请求体 | 任意类型，任意嵌套深度 |
| `path:"name"` | 路由参数 | 仅顶层标量字段 |
| `query:"name"` | URL 查询参数 | 仅顶层标量字段 |
| `validate:"rule"` | 校验规则 | 任意字段 |

- `path`/`query` 参数值为空时，不会覆盖数值/布尔/浮点字段。
- 内置校验规则：`required`、`min`、`max`、`email`、`url` 等，详见 [validator 文档](https://pkg.go.dev/github.com/go-playground/validator/v10)。
- 已注册自定义校验器 `phone`，用于校验中国大陆手机号（格式：`1[3-9]xxxxxxxxx`）。

## API

- `Parse[T any](c *gin.Context, obj *T) error`：将请求数据解析到 `obj`。
