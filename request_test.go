package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Index struct {
	Page int `form:"page" validate:"required,numeric,min=1"`
	Size int `form:"size" validate:"numeric,min=1,max=100"`
}

func (i *Index) CheckAndParseParams(c *gin.Context) error {
	return Parse(c, i)
}

func TestCheckAndParseParams_Success(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		params, err := CheckAndParseParams[Index](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data := params.Data()
		c.JSON(http.StatusOK, gin.H{"page": data.Page, "size": data.Size})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?page=2&size=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"page":2`)
	assert.Contains(t, w.Body.String(), `"size":10`)
}

func TestCheckAndParseParams_MissingRequired(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		params, err := CheckAndParseParams[Index](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, params.Data())
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?size=10", nil) // missing required page
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestCheckAndParseParams_SizeOutOfRange(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		params, err := CheckAndParseParams[Index](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, params.Data())
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?page=1&size=200", nil) // size > 100
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
