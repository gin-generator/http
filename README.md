# gin-generator/http

English | [中文](./README.zh-CN.md)

A lightweight helper for Gin that binds request body data and parses `path`/`query` parameters into a single struct.

## Features

- Uses `c.ShouldBind` to bind JSON/form data first, supporting nested structs, maps, and slices via `json` tags.
- Supports custom struct tags on top-level scalar fields:
  - `path:"name"` for route parameters
  - `query:"name"` for query parameters
- Scalar types supported for `path`/`query`:
  - `string`
  - signed integers (`int`, `int8`, `int16`, `int32`, `int64`)
  - unsigned integers (`uint`, `uint8`, `uint16`, `uint32`, `uint64`)
  - `bool`
  - floating numbers (`float32`, `float64`)

## Installation

```bash
go get github.com/gin-generator/http
```

## Usage

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
    Page    int               `query:"page"`
    Name    string            `json:"name"`
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

## How it works

1. `Parse` calls `c.ShouldBind(obj)` first to bind all `json`/`form` tagged fields, including nested structs, maps, and slices at any depth.
2. It then iterates through the top-level struct fields via reflection.
3. For each scalar field, it checks `path` and `query` tags and converts string values into the field type.
4. Non-scalar fields (`struct`, `map`, `slice`, etc.) are skipped — they are already populated by `ShouldBind`.

## Tag behavior

| Tag | Source | Applies to |
|---|---|---|
| `json:"name"` | request body | any type, any nesting depth |
| `path:"name"` | route parameter | top-level scalar fields only |
| `query:"name"` | URL query string | top-level scalar fields only |

- If a `path`/`query` value is empty, numeric/bool/float fields are not overwritten.

## API

- `Parse[T any](c *gin.Context, obj *T) error` — parse request data into `obj`.
