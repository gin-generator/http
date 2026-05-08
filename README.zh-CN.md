# gin-generator/http

[English](./README.md) | 中文

一个用于 Gin 的轻量请求解析辅助库：先绑定请求体数据，再把 `path`/`query` 参数解析到同一个结构体中。

## 功能特性

- 先通过 `c.ShouldBind` 绑定 JSON/form 数据。
- 支持自定义结构体标签：
  - `path:"name"`：路由参数
  - `query:"name"`：查询参数
- 支持以下字段类型解析：
  - `string`
  - 有符号整数（`int`、`int8`、`int16`、`int32`、`int64`）
  - 无符号整数（`uint`、`uint8`、`uint16`、`uint32`、`uint64`）
  - `bool`
  - 浮点数（`float32`、`float64`）

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

type GetUserRequest struct {
    ID      int    `path:"id"`
    Keyword string `query:"keyword"`
    Active  bool   `query:"active"`
    Name    string `json:"name"`
}

func main() {
    r := gin.Default()

    r.POST("/users/:id", func(c *gin.Context) {
        var req GetUserRequest
        if err := _http.Parse(c, &req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "id":      req.ID,
            "keyword": req.Keyword,
            "active":  req.Active,
            "name":    req.Name,
        })
    })

    _ = r.Run(":8080")
}
```

## 工作原理

1. `Parse` 先调用 `c.ShouldBind(obj)`，绑定请求体或表单数据。
2. 使用反射遍历结构体字段。
3. 读取每个字段的 `path`/`query` 标签，把字符串参数转换为对应字段类型。
4. 转换失败时直接返回错误。

## 标签行为说明

- 字段没有 `path`/`query` 标签时，`ShouldBind` 后不会被二次覆盖。
- 标签存在但参数值为空时，不会覆盖数值/布尔/浮点字段。
- 遇到不支持的字段类型会返回错误。

## API

- `Parse[T any](c *gin.Context, obj *T) error`：将请求数据解析到 `obj`。
