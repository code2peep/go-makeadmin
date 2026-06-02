package routers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-makeadmin/core"
	"go-makeadmin/core/response"
)

func TestInitRoutersDefaultEmpty(t *testing.T) {
	t.Setenv(EnableDemoModuleEnv, "")

	if got := InitRouters(); len(got) != 0 {
		t.Fatalf("expected no demo runtime routers by default, got %d", len(got))
	}
}

func TestDemoRuntimeRouteRequiresToken(t *testing.T) {
	t.Setenv(EnableDemoModuleEnv, "1")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	group := router.Group("/api")
	for _, item := range InitRouters() {
		core.RegisterGroup(group, item)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/article/list", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 wrapper response, got %d", res.Code)
	}

	var body response.Response
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != response.TokenEmpty.Code() {
		t.Fatalf("expected token empty response, got code %d", body.Code)
	}
}
