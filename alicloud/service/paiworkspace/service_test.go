package paiworkspace

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// paiWorkspaceTestClient stands up a credential-free AliyunClient whose
// AIWorkSpace ROA transport is pointed at an in-process httptest server,
// mirroring the rdsAccountTestClient construction. It exercises the real ROA
// request/signing/parse path without touching the cloud, so the Create and
// Read wire layer is covered without placing a chargeable order.
func paiWorkspaceTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("aiworkspace", endpoint)
	return client
}

func writeJSON(t *testing.T, w http.ResponseWriter, v map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Error(err)
	}
}

func TestUnitBuildCreateBody(t *testing.T) {
	// Minimal valid order: ProductCode is always present, AutoRenew/AutoPay
	// default to false, no optional order fields.
	got := buildCreateBody(productCreateInput{ProductCode: "PAI_share"})
	product := got["Products"].([]interface{})[0].(map[string]interface{})
	if product["ProductCode"] != "PAI_share" {
		t.Fatalf("ProductCode = %v, want PAI_share", product["ProductCode"])
	}
	if product["AutoRenew"] != false {
		t.Fatalf("AutoRenew = %v, want false", product["AutoRenew"])
	}
	if got["AutoPay"] != false {
		t.Fatalf("AutoPay = %v, want false", got["AutoPay"])
	}
	if _, ok := product["ChargeType"]; ok {
		t.Fatalf("ChargeType should be absent when empty, got %v", product["ChargeType"])
	}
	if _, ok := product["InstanceProperties"]; ok {
		t.Fatalf("InstanceProperties should be absent when empty")
	}

	// Full order: all optional fields are sent; instance_properties preserve
	// only non-empty code/name/value entries.
	got = buildCreateBody(productCreateInput{
		ProductCode: "PAI_share", ChargeType: "POSTPAY", PricingCycle: "Hour",
		Duration: 1, AutoRenew: true, OrderType: "BUY", AutoPay: true,
		InstanceProperties: []instancePropertyInput{
			{Code: "region", Name: "Region", Value: "cn-hangzhou"},
			{Code: "", Name: "", Value: ""}, // dropped: all empty
		},
	})
	product = got["Products"].([]interface{})[0].(map[string]interface{})
	if product["ChargeType"] != "POSTPAY" || product["PricingCycle"] != "Hour" ||
		product["Duration"] != int64(1) || product["OrderType"] != "BUY" ||
		product["AutoRenew"] != true || got["AutoPay"] != true {
		t.Fatalf("full order product = %#v", product)
	}
	ips := product["InstanceProperties"].([]interface{})
	if len(ips) != 1 {
		t.Fatalf("InstanceProperties len = %d, want 1 (empty entry dropped)", len(ips))
	}
	ip := ips[0].(map[string]interface{})
	if ip["Code"] != "region" || ip["Name"] != "Region" || ip["Value"] != "cn-hangzhou" {
		t.Fatalf("instance property = %#v, want region/Region/cn-hangzhou", ip)
	}
}

func TestUnitFindProduct(t *testing.T) {
	products := []interface{}{
		map[string]interface{}{"ProductId": "pid-1", "ProductCode": "PAI_share"},
		map[string]interface{}{"ProductId": "pid-2", "ProductCode": "other"},
		"not-a-map",
	}
	if got := findProduct(products, "pid-1"); got["ProductCode"] != "PAI_share" {
		t.Fatalf("findProduct(pid-1) = %#v", got)
	}
	if got := findProduct(products, "missing"); got != nil {
		t.Fatalf("findProduct(missing) = %#v, want nil", got)
	}
	if got := findProduct(products, "pid-2"); got["ProductCode"] != "other" {
		t.Fatalf("findProduct(pid-2) = %#v", got)
	}
}

func TestUnitApplyProductObjectPreservesCreateOnlyFields(t *testing.T) {
	// Create-only order attributes and instance_properties are not returned by
	// ListProducts; applyProductObject must leave them untouched so refresh
	// does not drop configured values.
	state := &productResourceModel{
		ChargeType:   types.StringValue("POSTPAY"),
		PricingCycle: types.StringValue("Hour"),
		Duration:     types.Int64Value(1),
		OrderType:    types.StringValue("BUY"),
		AutoRenew:    types.BoolValue(true),
		AutoPay:      types.BoolValue(true),
		ProductCode:  types.StringValue("PAI_share"),
	}
	obj := map[string]interface{}{
		"ProductId":               "pid-9",
		"ProductCode":             "PAI_share",
		"IsPurchased":             true,
		"HasPermissionToPurchase": true,
		"PurchaseUrl":             "https://example.com/buy",
		"RegionId":                "cn-hangzhou",
	}
	applyProductObject(state, obj, "cn-hangzhou")
	if state.ProductId.ValueString() != "pid-9" || state.Id.ValueString() != "pid-9" {
		t.Fatalf("ProductId/Id = %q/%q, want pid-9", state.ProductId, state.Id)
	}
	if !state.IsPurchased.ValueBool() || !state.HasPermissionToPurchase.ValueBool() {
		t.Fatalf("catalog bools not mapped: %+v", state)
	}
	if state.PurchaseUrl.ValueString() != "https://example.com/buy" {
		t.Fatalf("PurchaseUrl = %q", state.PurchaseUrl)
	}
	// Create-only fields preserved from prior state, not overwritten by the
	// (absent) catalog values.
	if state.ChargeType.ValueString() != "POSTPAY" || state.PricingCycle.ValueString() != "Hour" ||
		state.OrderType.ValueString() != "BUY" || state.Duration.ValueInt64() != 1 ||
		!state.AutoRenew.ValueBool() || !state.AutoPay.ValueBool() {
		t.Fatalf("create-only fields clobbered by applyProductObject: %+v", state)
	}
	// RegionId falls back to the client region when the catalog omits it.
	applyProductObject(state, map[string]interface{}{"ProductId": "pid-9"}, "cn-hangzhou")
	if state.RegionId.ValueString() != "cn-hangzhou" {
		t.Fatalf("RegionId fallback = %q, want cn-hangzhou", state.RegionId)
	}
}

func TestUnitCreateProductMockSuccess(t *testing.T) {
	var seenPath bool
	client := paiWorkspaceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "productorders") {
			seenPath = true
			writeJSON(t, w, map[string]interface{}{"ProductIds": []string{"pid-mock-1"}})
			return
		}
		http.NotFound(w, r)
	})
	id, err := createProduct(client, buildCreateBody(productCreateInput{ProductCode: "PAI_share"}))
	if err != nil {
		t.Fatalf("createProduct error: %v", err)
	}
	if id != "pid-mock-1" {
		t.Fatalf("createProduct id = %q, want pid-mock-1", id)
	}
	if !seenPath {
		t.Fatal("CreateProductOrders path was not hit")
	}
}

func TestUnitCreateProductMockFailure(t *testing.T) {
	// Failure path on the create API: the server returns an error code body,
	// createProduct surfaces it as an error and never returns a ProductId. No
	// order is placed on the mock, mirroring a real rejected order (no charge).
	client := paiWorkspaceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "productorders") {
			writeJSON(t, w, map[string]interface{}{
				"Code": "InvalidProduct.ProductCode", "Message": "product code not found",
				"RequestId": "req-1", "HostId": "h-1",
			})
			return
		}
		http.NotFound(w, r)
	})
	id, err := createProduct(client, buildCreateBody(productCreateInput{ProductCode: "INVALID"}))
	if err == nil {
		t.Fatalf("createProduct expected error, got id=%q", id)
	}
	if id != "" {
		t.Fatalf("createProduct id = %q on failure, want empty", id)
	}
	if !strings.Contains(err.Error(), "InvalidProduct.ProductCode") {
		t.Fatalf("error does not carry the API code: %v", err)
	}
}

func TestUnitCreateProductMockEmptyIds(t *testing.T) {
	client := paiWorkspaceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]interface{}{"ProductIds": []interface{}{}})
	})
	_, err := createProduct(client, buildCreateBody(productCreateInput{ProductCode: "PAI_share"}))
	if err == nil {
		t.Fatal("createProduct expected error for empty ProductIds")
	}
}

func TestUnitReadProductMockFound(t *testing.T) {
	client := paiWorkspaceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "products") {
			if c := r.URL.Query().Get("ProductCodes"); c != "PAI_share" {
				t.Errorf("ProductCodes query = %q, want PAI_share", c)
			}
			writeJSON(t, w, map[string]interface{}{
				"Products": []map[string]interface{}{
					{"ProductId": "pid-mock-1", "ProductCode": "PAI_share", "IsPurchased": true},
				},
			})
			return
		}
		http.NotFound(w, r)
	})
	obj, err := readProduct(client, "pid-mock-1", "PAI_share")
	if err != nil {
		t.Fatalf("readProduct error: %v", err)
	}
	if obj == nil || obj["ProductId"] != "pid-mock-1" {
		t.Fatalf("readProduct obj = %#v", obj)
	}
}

func TestUnitReadProductMockNotFound(t *testing.T) {
	client := paiWorkspaceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]interface{}{
			"Products": []map[string]interface{}{
				{"ProductId": "pid-other", "ProductCode": "PAI_share"},
			},
		})
	})
	obj, err := readProduct(client, "pid-missing", "PAI_share")
	if err != nil {
		t.Fatalf("readProduct error: %v", err)
	}
	if obj != nil {
		t.Fatalf("readProduct obj = %#v, want nil (not found)", obj)
	}
}

func TestUnitDeleteIsNoOp(t *testing.T) {
	// Delete has no API: orders are irreversible. The resource Delete method
	// must not call the cloud and must not emit diagnostics.
	r := &productResource{}
	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), resource.DeleteRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete produced diagnostics: %v", resp.Diagnostics.Errors())
	}
}
