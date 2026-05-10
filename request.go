package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CheckAndParseParams(c *gin.Context) error
}

// Request can be used to verify at compile time that *T implements Handler.
type Request[T any] interface {
	*T
	Handler
}

type RequestType[T any] struct {
	data T
}

func NewRequestType[T any](data T) RequestType[T] {
	return RequestType[T]{data: data}
}

func (r RequestType[T]) Data() T {
	return r.data
}

func CheckAndParseParams[T any](c *gin.Context) (RequestType[T], error) {
	var v T
	h, ok := any(&v).(Handler)
	if !ok {
		panic(fmt.Sprintf("*%T does not implement Handler", v))
	}
	if err := h.CheckAndParseParams(c); err != nil {
		return RequestType[T]{}, err
	}
	return RequestType[T]{v}, nil
}
