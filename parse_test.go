package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	UserID  int    `uri:"user_id"`
	Page    int    `form:"page" validate:"required,min=1"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"phone"`
	Email   string `json:"email" validate:"omitempty,email"`
	Age     int    `json:"age" validate:"omitempty,min=1,max=150"`
	Enabled bool   `form:"enabled"`
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestParse_Success(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	body := `{"name":"Alice","phone":"13800138000","email":"alice@example.com","age":25}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=2&enabled=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestParse_ValidationError_RequiredField(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Missing required "name" field
	body := `{"phone":"13800138000"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_MinValue(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// page=0 violates min=1
	body := `{"name":"Alice","phone":"13800138000"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=0", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_CustomPhone(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Invalid phone format
	body := `{"name":"Alice","phone":"12345"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_Email(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Invalid email format
	body := `{"name":"Alice","phone":"13800138000","email":"invalid-email"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_ValidationError_AgeRange(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// age=200 violates max=150
	body := `{"name":"Alice","phone":"13800138000","age":200}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders?page=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

type UpdateUserRequest struct {
	UserID int    `uri:"user_id" validate:"required"`
	OrgID  int    `uri:"org_id" validate:"required"`
	Notify bool   `form:"notify"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
}

func TestParse_PUT_UriAndQueryAndJSON(t *testing.T) {
	router := setupRouter()
	router.PUT("/orgs/:org_id/users/:user_id", func(c *gin.Context) {
		var req UpdateUserRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	body := `{"name":"Alice","email":"alice@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/orgs/7/users/123?notify=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"UserID":123`)
	assert.Contains(t, w.Body.String(), `"OrgID":7`)
	assert.Contains(t, w.Body.String(), `"Notify":true`)
	assert.Contains(t, w.Body.String(), `"Name":"Alice"`)
	assert.Contains(t, w.Body.String(), `"Email":"alice@example.com"`)
}

func TestParse_PUT_MissingUriParam(t *testing.T) {
	router := setupRouter()
	// Route without org_id so the uri tag can't be filled.
	router.PUT("/users/:user_id", func(c *gin.Context) {
		var req UpdateUserRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	body := `{"name":"Alice","email":"alice@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_PUT_InvalidJSONBody(t *testing.T) {
	router := setupRouter()
	router.PUT("/orgs/:org_id/users/:user_id", func(c *gin.Context) {
		var req UpdateUserRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Invalid email triggers validator failure even though uri/query are fine.
	body := `{"name":"Alice","email":"not-an-email"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/orgs/7/users/123?notify=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestParse_MissingQueryParam(t *testing.T) {
	router := setupRouter()
	router.POST("/users/:user_id/orders", func(c *gin.Context) {
		var req TestRequest
		if err := Parse(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// Missing required query param "page"
	body := `{"name":"Alice","phone":"13800138000"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
