# gin-generator/http

[English](./README.md) | 中文

一个用于 Gin 的轻量请求解析辅助库：将 URI / query / form / JSON 数据绑定到同一个结构体，自动校验，并支持可选的解析后钩子。

## 功能特性

- 使用 Gin 内置的绑定器解析全部请求来源：
  - `uri:"name"`：路由参数（`ShouldBindUri`）
  - `form:"name"`：查询字符串与表单字段（`ShouldBindQuery` + `ShouldBind`）
  - `json:"name"`：JSON 请求体（`ShouldBind`）
- 使用 `validate` 标签自动校验，基于 [go-playground/validator](https://github.com/go-playground/validator)。
- 可选的 `AfterParser` 钩子，便于在绑定与校验后执行跨字段检查或自定义逻辑。
- 泛型入口 `CheckAndParseParams[T]` 返回类型化的 `RequestType[T]` 包装。

## 安装

```bash
go get github.com/gin-generator/http
```

## 使用示例

```go
package main

import (
    "fmt"
    "net/http"

    _http "github.com/gin-generator/http"
    "github.com/gin-gonic/gin"
)

type Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type CreateOrderRequest struct {
    UserID  int               `uri:"user_id" validate:"required"`
    Page    int               `form:"page" validate:"required,min=1"`
    Name    string            `json:"name" validate:"required"`
    Phone   string            `json:"phone" validate:"phone"`
    Tags    []string          `json:"tags"`
    Meta    map[string]string `json:"meta"`
    Address Address           `json:"address"`
    Items   []Address         `json:"items"`
}

// 可选：实现 AfterParser，用于跨字段校验。
func (r *CreateOrderRequest) AfterParse(c *gin.Context) error {
    if len(r.Items) == 0 {
        return fmt.Errorf("items must not be empty")
    }
    return nil
}

func main() {
    r := gin.Default()

    r.POST("/users/:user_id/orders", func(c *gin.Context) {
        params, err := _http.CheckAndParseParams[CreateOrderRequest](c)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, params.Data())
    })

    _ = r.Run(":8080")
}
```

## 工作原理

1. `Parse` 分三步独立绑定，单一来源缺失不会影响其它来源：
   - `c.ShouldBindUri(obj)`：路由参数
   - `c.ShouldBindQuery(obj)`：查询字符串
   - `c.ShouldBind(obj)`：请求体（按 Content-Type 自动选择 JSON / form 等）
2. 最后调用 `ValidateStruct` 对整个结构体执行一次校验。
3. `CheckAndParseParams[T]` 随后判断 `*T` 是否实现了 `AfterParser`，若实现则执行 `AfterParse(c)`，用于自定义的后置校验逻辑。
4. 成功后将解析结果包装在 `RequestType[T]` 中，使用 `.Data()` 获取值。

## 标签行为说明

| 标签 | 数据来源 | 绑定方式 |
|---|---|---|
| `uri:"name"` | 路由参数 | `ShouldBindUri` |
| `form:"name"` | 查询字符串 / 表单请求体 | `ShouldBindQuery` + `ShouldBind` |
| `json:"name"` | JSON 请求体 | `ShouldBind` |
| `validate:"rule"` | 校验规则 | `ValidateStruct` |

- 内置校验规则：`required`、`min`、`max`、`email`、`url` 等，详见 [validator 文档](https://pkg.go.dev/github.com/go-playground/validator/v10)。
- 已注册自定义校验器 `phone`，用于校验中国大陆手机号（格式：`1[3-9]xxxxxxxxx`）。

## API

- `Parse[T any](c *gin.Context, obj *T) error`：将 uri/query/body 绑定到 `obj` 并执行校验。
- `CheckAndParseParams[T any](c *gin.Context) (RequestType[T], error)`：`Parse` + 可选的 `AfterParser` 钩子，返回类型化包装。
- `RequestType[T].Data() T`：获取解析后的值。
- `AfterParser`：可选接口，在 `*T` 上实现 `AfterParse(c *gin.Context) error`，用于绑定与校验后的自定义逻辑。
- `ValidateStruct(s interface{}) error`：使用共享的 validator 实例校验结构体。
