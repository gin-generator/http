# gin-generator/http

English | [中文](./README.zh-CN.md)

A lightweight helper for Gin that binds URI / query / form / JSON data into a single struct, validates it, and optionally runs a post-parse hook.

## Features

- Binds all request sources via Gin's built-in binders:
  - `uri:"name"` — route parameters (`ShouldBindUri`)
  - `form:"name"` — query string and form fields (`ShouldBindQuery` + `ShouldBind`)
  - `json:"name"` — JSON body (`ShouldBind`)
- Automatic validation using `validate` tags powered by [go-playground/validator](https://github.com/go-playground/validator).
- Optional `AfterParser` hook for cross-field checks or custom logic after parsing and validation.
- Generic `CheckAndParseParams[T]` returns a typed `RequestType[T]` wrapper.

## Installation

```bash
go get github.com/gin-generator/http
```

## Usage

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

// Optional: implement AfterParser for cross-field validation.
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

## How it works

1. `Parse` binds in three steps, each independent so missing sources don't fail the whole call:
   - `c.ShouldBindUri(obj)` — route parameters
   - `c.ShouldBindQuery(obj)` — query string
   - `c.ShouldBind(obj)` — body (JSON / form, content-type aware)
2. The struct is validated once at the end via `ValidateStruct`.
3. `CheckAndParseParams[T]` then checks if `*T` implements `AfterParser`; if so, `AfterParse(c)` runs for custom post-validation logic.
4. On success, the parsed value is wrapped in `RequestType[T]`; access it with `.Data()`.

## Tag behavior

| Tag | Source | Bound by |
|---|---|---|
| `uri:"name"` | route parameter | `ShouldBindUri` |
| `form:"name"` | query string / form body | `ShouldBindQuery` + `ShouldBind` |
| `json:"name"` | JSON request body | `ShouldBind` |
| `validate:"rule"` | validation rule | `ValidateStruct` |

- Built-in validation rules: `required`, `min`, `max`, `email`, `url`, etc. See [validator docs](https://pkg.go.dev/github.com/go-playground/validator/v10).
- Custom validator `phone` is registered for Chinese mobile phone numbers (format: `1[3-9]xxxxxxxxx`).

## API

- `Parse[T any](c *gin.Context, obj *T) error` — bind uri/query/body into `obj` and validate it.
- `CheckAndParseParams[T any](c *gin.Context) (RequestType[T], error)` — `Parse` + optional `AfterParser` hook, returns a typed wrapper.
- `RequestType[T].Data() T` — retrieve the parsed value.
- `AfterParser` — optional interface; implement `AfterParse(c *gin.Context) error` on `*T` to run custom logic after binding and validation.
- `ValidateStruct(s interface{}) error` — validate a struct using the shared validator instance.
