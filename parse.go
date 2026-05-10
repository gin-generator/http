package http

import (
	"reflect"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	tagPath  = "path"
	tagQuery = "query"
)

type fieldMeta struct {
	index    int
	kind     reflect.Kind
	pathTag  string
	queryTag string
}

var metaCache sync.Map

// getStructMeta builds and caches field metadata for path/query injection.
// Only scalar fields with a path or query tag are included.
func getStructMeta(t reflect.Type) []fieldMeta {
	if v, ok := metaCache.Load(t); ok {
		return v.([]fieldMeta)
	}
	fields := make([]fieldMeta, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		pathTag := f.Tag.Get(tagPath)
		queryTag := f.Tag.Get(tagQuery)
		if pathTag == "" && queryTag == "" {
			continue
		}
		kind := f.Type.Kind()
		switch kind {
		case reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Bool,
			reflect.Float32, reflect.Float64:
			fields = append(fields, fieldMeta{
				index:    i,
				kind:     kind,
				pathTag:  pathTag,
				queryTag: queryTag,
			})
		default:
			panic("unhandled default case")
		}
	}
	metaCache.Store(t, fields)
	return fields
}

// Parse binds request data (JSON/form/query) and injects path/query params.
// JSON body and form fields are handled by ShouldBind via their respective tags.
// path/query tags are supported on top-level scalar fields only.
// The struct is validated after all fields are populated.
func Parse[T any](c *gin.Context, obj *T) error {
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	val := reflect.ValueOf(obj).Elem()
	for _, meta := range getStructMeta(val.Type()) {
		if err := setField(c, val.Field(meta.index), meta); err != nil {
			return err
		}
	}
	return ValidateStruct(obj)
}

func setField(c *gin.Context, f reflect.Value, meta fieldMeta) error {
	var raw string
	if meta.pathTag != "" {
		raw = c.Param(meta.pathTag)
	} else if meta.queryTag != "" {
		raw = c.Query(meta.queryTag)
	}
	if raw == "" {
		return nil
	}
	switch meta.kind {
	case reflect.String:
		f.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		f.SetUint(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		f.SetFloat(v)
	default:
		panic("unhandled default case")
	}
	return nil
}
