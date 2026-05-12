package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Index struct {
	Page int `form:"page" validate:"required,min=1"`
	Size int `form:"size" validate:"min=1,max=100"`
}

type IndexWithHook struct {
	Page int `form:"page" validate:"required,min=1"`
	Size int `form:"size" validate:"min=1,max=100"`
}

func (i *IndexWithHook) AfterParse(c *gin.Context) error {
	if i.Page*i.Size > 1000 {
		return fmt.Errorf("page*size too large")
	}
	return nil
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

// --- AfterParser hook tests ---

func TestCheckAndParseParams_AfterParse_Success(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		params, err := CheckAndParseParams[IndexWithHook](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data := params.Data()
		c.JSON(http.StatusOK, gin.H{"page": data.Page, "size": data.Size})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?page=2&size=10", nil) // 2*10=20 <= 1000
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"page":2`)
}

func TestCheckAndParseParams_AfterParse_Rejected(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		params, err := CheckAndParseParams[IndexWithHook](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, params.Data())
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?page=100&size=100", nil) // 100*100=10000 > 1000
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "page*size too large")
}

// NoHook struct: no AfterParser, hook branch should be silently skipped.
type NoHookIndex struct {
	Page int `form:"page" validate:"required,min=1"`
}

func TestCheckAndParseParams_NoHook(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(c *gin.Context) {
		params, err := CheckAndParseParams[NoHookIndex](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"page": params.Data().Page})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?page=5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"page":5`)
}
