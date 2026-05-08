package http

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	tagPath  = "path"
	tagQuery = "query"
)

// Parse binds request data (JSON/form/query) and extracts path/query parameters into the struct.
// It supports basic types: string, int, uint, bool, float.
//
// Example:
//
//	type Request struct {
//	    ID   int    `path:"id"`
//	    Name string `json:"name"`
//	    Page int    `query:"page"`
//	}
//	var req Request
//	if err := Parse(c, &req); err != nil {
//	    // handle error
//	}
func Parse[T any](c *gin.Context, obj *T) error {
	// Bind JSON/form data first
	if err := c.ShouldBind(obj); err != nil {
		return err
	}

	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := typ.Field(i)
		tag := fieldType.Tag

		if err := parseField(c, &field, tag); err != nil {
			return err
		}
	}

	return nil
}

// parseField parses a single field based on its type and tags.
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
		return errors.New("unsupported field type: " + field.Kind().String())
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
