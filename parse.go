package http

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	tagPath  = "path"
	tagQuery = "query"
)

// Parse binds request data (JSON/form/query) and extracts path/query parameters into the struct.
// JSON body binding (including nested structs, maps, slices) is handled by ShouldBind.
// path/query tags are only supported on top-level scalar fields (string, int, uint, bool, float).
//
// Example:
//
//	type Request struct {
//	    ID     int               `path:"id"`
//	    Page   int               `query:"page"`
//	    Name   string            `json:"name"`
//	    Tags   []string          `json:"tags"`
//	    Meta   map[string]string `json:"meta"`
//	}
//	var req Request
//	if err := Parse(c, &req); err != nil {
//	    // handle error
//	}
func Parse[T any](c *gin.Context, obj *T) error {
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if !f.CanSet() {
			continue
		}
		if err := parseField(c, &f, typ.Field(i).Tag); err != nil {
			return err
		}
	}
	return nil
}

// parseField handles path/query tag injection for scalar fields.
// struct/map/slice fields are json-only and already populated by ShouldBind — skip them.
func parseField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) error {
	switch field.Kind() {
	case reflect.String:
		return parseStringField(c, field, tag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseIntField(c, field, tag)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseUintField(c, field, tag)
	case reflect.Bool:
		return parseBoolField(c, field, tag)
	case reflect.Float32, reflect.Float64:
		return parseFloatField(c, field, tag)
	default:
		// map, slice, interface, etc. — already handled by ShouldBind, skip.
		return nil
	}
}

func parseStringField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) error {
	if pathTag, ok := tag.Lookup(tagPath); ok {
		field.SetString(c.Param(pathTag))
	} else if queryTag, okk := tag.Lookup(tagQuery); okk {
		field.SetString(c.Query(queryTag))
	}
	return nil
}

func parseIntField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) error {
	var value string
	if pathTag, ok := tag.Lookup(tagPath); ok {
		value = c.Param(pathTag)
	} else if queryTag, okk := tag.Lookup(tagQuery); okk {
		value = c.Query(queryTag)
	}

	if value != "" {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)
	}
	return nil
}

func parseUintField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) error {
	var value string
	if pathTag, ok := tag.Lookup(tagPath); ok {
		value = c.Param(pathTag)
	} else if queryTag, okk := tag.Lookup(tagQuery); okk {
		value = c.Query(queryTag)
	}

	if value != "" {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(v)
	}
	return nil
}

func parseBoolField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) error {
	var value string
	if pathTag, ok := tag.Lookup(tagPath); ok {
		value = c.Param(pathTag)
	} else if queryTag, okk := tag.Lookup(tagQuery); okk {
		value = c.Query(queryTag)
	}

	if value != "" {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)
	}
	return nil
}

func parseFloatField(c *gin.Context, field *reflect.Value, tag reflect.StructTag) error {
	var value string
	if pathTag, ok := tag.Lookup(tagPath); ok {
		value = c.Param(pathTag)
	} else if queryTag, okk := tag.Lookup(tagQuery); okk {
		value = c.Query(queryTag)
	}

	if value != "" {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)
	}
	return nil
}
