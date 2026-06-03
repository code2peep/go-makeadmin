package gen

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-makeadmin/core/response"
	gensvc "go-makeadmin/generator/service/gen"
)

type moduleRegistryResponseItem struct {
	Name            string                     `json:"name"`
	Module          string                     `json:"module"`
	Manifest        string                     `json:"manifest"`
	Table           string                     `json:"table"`
	Runtime         string                     `json:"runtime"`
	Entry           string                     `json:"entry"`
	ManifestStatus  string                     `json:"manifestStatus"`
	ManifestMessage string                     `json:"manifestMessage"`
	ManifestChecks  []moduleRegistryCheckEntry `json:"manifestChecks"`
}

type moduleRegistryCheckEntry struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestListModuleRegistryRouteDefaultResponse(t *testing.T) {
	t.Setenv(gensvc.EnableBrokenModuleRegistryFixtureEnv, "")

	items := moduleRegistryRouteResponse(t)
	if len(items) != 1 {
		t.Fatalf("registry route items = %d, want 1", len(items))
	}
	item := items[0]
	if item.Module != "article" || item.Manifest != "examples/demo/manifest.json" || item.Entry != "/demo/article" {
		t.Fatalf("unexpected route item: %+v", item)
	}
	if item.ManifestStatus != "passed" || item.ManifestMessage == "" || len(item.ManifestChecks) < 6 {
		t.Fatalf("unexpected route manifest checks: %+v", item)
	}
}

func TestListModuleRegistryRouteBrokenFixtureResponse(t *testing.T) {
	t.Setenv(gensvc.EnableBrokenModuleRegistryFixtureEnv, "1")

	items := moduleRegistryRouteResponse(t)
	if len(items) != 2 {
		t.Fatalf("registry route items = %d, want 2", len(items))
	}
	demo := findModuleRegistryRouteItem(items, "article")
	if demo == nil || demo.ManifestStatus != "passed" {
		t.Fatalf("demo route item must still pass: %+v", demo)
	}
	broken := findModuleRegistryRouteItem(items, "broken_fixture")
	if broken == nil {
		t.Fatalf("broken fixture route item not found: %+v", items)
	}
	if broken.ManifestStatus != "failed" || broken.ManifestMessage == "" {
		t.Fatalf("broken fixture route item must fail: %+v", broken)
	}
	if len(broken.ManifestChecks) != 1 || broken.ManifestChecks[0].Status != "failed" {
		t.Fatalf("unexpected broken fixture route checks: %+v", broken.ManifestChecks)
	}
}

func moduleRegistryRouteResponse(t *testing.T) []moduleRegistryResponseItem {
	t.Helper()
	gin.SetMode(gin.TestMode)

	handler := newGenHandler(gensvc.NewGenerateService(nil))
	req := httptest.NewRequest(http.MethodGet, "/api/gen/moduleRegistry", nil)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = req

	handler.listModuleRegistry(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 wrapper response, got %d", recorder.Code)
	}
	var body struct {
		Code int                          `json:"code"`
		Data []moduleRegistryResponseItem `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode route response: %v", err)
	}
	if body.Code != response.Success.Code() {
		t.Fatalf("expected success code, got %d", body.Code)
	}
	return body.Data
}

func findModuleRegistryRouteItem(items []moduleRegistryResponseItem, module string) *moduleRegistryResponseItem {
	for index := range items {
		if items[index].Module == module {
			return &items[index]
		}
	}
	return nil
}
