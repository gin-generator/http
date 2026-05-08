# gin-generator/http

English | [中文](./README.zh-CN.md)

A lightweight helper for Gin that binds request body data and parses `path`/`query` parameters into a single struct.

## Features

- Uses `c.ShouldBind` to bind JSON/form data first.
- Supports custom struct tags:
  - `path:"name"` for route parameters
  - `query:"name"` for query parameters
- Parses and sets these field kinds:
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

## How it works

1. `Parse` calls `c.ShouldBind(obj)` first to bind request body/form data.
2. It then iterates through struct fields via reflection.
3. For each field, it checks `path` and `query` tags and converts string values into the field type.
4. Returns conversion errors directly when parsing fails.

## Tag behavior

- If both `path` and `query` are absent on a field, `Parse` leaves that field unchanged after `ShouldBind`.
- If a tagged value is empty, numeric/bool/float fields are not overwritten.
- Unsupported field kinds return an error.

## API

- `Parse[T any](c *gin.Context, obj *T) error` — parse request data into `obj`.
