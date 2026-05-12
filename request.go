package http

import (
	"github.com/gin-gonic/gin"
)

// AfterParser is an optional hook. Implement it on *T to run custom logic
// (e.g. cross-field checks) after binding and tag validation.
type AfterParser interface {
	AfterParse(c *gin.Context) (err error)
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

// Parse binds uri/query/form/json into obj and validates it.
// Use uri tag for route parameters, form tag for query string and form fields, json tag for JSON body.
func Parse[T any](c *gin.Context, obj *T) error {
	// Each source is bound independently; validation runs once at the end.
	_ = c.ShouldBindUri(obj)
	_ = c.ShouldBindQuery(obj)
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	return ValidateStruct(obj)
}

func CheckAndParseParams[T any](c *gin.Context) (RequestType[T], error) {
	var v T
	if err := Parse(c, &v); err != nil {
		return RequestType[T]{}, err
	}
	if h, ok := any(&v).(AfterParser); ok {
		if err := h.AfterParse(c); err != nil {
			return RequestType[T]{}, err
		}
	}
	return RequestType[T]{v}, nil
}
