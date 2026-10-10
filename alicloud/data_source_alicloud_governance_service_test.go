package alicloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudGovernanceServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_governance_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudGovernanceServiceDataSourceOpen,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudGovernanceServiceDataSourceDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "id", "GovernanceServiceHasNotBeenOpened"),
					resource.TestCheckResourceAttr(resourceId, "status", ""),
				),
			},
			{
				Config: testAccCheckAlicloudGovernanceServiceDataSourceOpen,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudGovernanceServiceDataSourceOpen = `
data "alicloud_governance_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudGovernanceServiceDataSourceDefault = `
data "alicloud_governance_service" "current" {
}
`

// testUnitGovernanceServiceStubServer returns a local stub for the governance
// RPC endpoint: it records every request's Action and form so the Read paths
// are locked against real network access while exercising the real
// AliyunClient.rpcRequest chain (endpoint resolution, signing, error format).
func testUnitGovernanceServiceStubServer(t *testing.T, statusCode int, body string, record func(action string, form map[string]string)) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		form := map[string]string{}
		for key := range r.Form {
			form[key] = r.Form.Get(key)
		}
		record(r.Form.Get("Action"), form)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return server
}

func testUnitGovernanceServiceMockMeta(t *testing.T, server *httptest.Server) *connectivity.AliyunClient {
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	host := server.Listener.Addr().String()
	t.Setenv("NO_PROXY", host)
	endpoints := new(sync.Map)
	endpoints.Store("governance", host)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	meta, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return meta
}

func TestUnitAliCloudGovernanceServiceRead(t *testing.T) {
	t.Run("EnableOnOpensService", func(t *testing.T) {
		var actions []string
		var capturedQuery map[string]string
		server := testUnitGovernanceServiceStubServer(t, http.StatusOK, `{"RequestId":"stub"}`, func(action string, form map[string]string) {
			actions = append(actions, action)
			capturedQuery = form
		})

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudGovernanceService().Schema, map[string]interface{}{
			"enable": "On",
		})
		if err := dataSourceAliCloudGovernanceServiceRead(d, testUnitGovernanceServiceMockMeta(t, server)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"OpenGovernanceService"}, actions)
		assert.NotEmpty(t, capturedQuery["RegionId"])
		assert.Equal(t, "GovernanceServiceHasBeenOpened", d.Id())
		assert.Equal(t, "Opened", d.Get("status"))
	})

	t.Run("EnableOffDoesNotCallOpen", func(t *testing.T) {
		var actions []string
		server := testUnitGovernanceServiceStubServer(t, http.StatusOK, `{"RequestId":"unused"}`, func(action string, _ map[string]string) {
			actions = append(actions, action)
		})

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudGovernanceService().Schema, map[string]interface{}{
			"enable": "Off",
		})
		if err := dataSourceAliCloudGovernanceServiceRead(d, testUnitGovernanceServiceMockMeta(t, server)); err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, actions)
		assert.Equal(t, "GovernanceServiceHasNotBeenOpened", d.Id())
		assert.Equal(t, "", d.Get("status"))
	})

	t.Run("EnableMissingDefaultsToOffNoCall", func(t *testing.T) {
		var actions []string
		server := testUnitGovernanceServiceStubServer(t, http.StatusOK, `{"RequestId":"unused"}`, func(action string, _ map[string]string) {
			actions = append(actions, action)
		})

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudGovernanceService().Schema, map[string]interface{}{})
		if err := dataSourceAliCloudGovernanceServiceRead(d, testUnitGovernanceServiceMockMeta(t, server)); err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, actions)
		assert.Equal(t, "GovernanceServiceHasNotBeenOpened", d.Id())
		assert.Equal(t, "", d.Get("status"))
	})

	t.Run("EnableOnAlreadyOpenedTolerated", func(t *testing.T) {
		var actions []string
		server := testUnitGovernanceServiceStubServer(t, http.StatusForbidden, `{"RequestId":"stub","HostId":"governance.cn-hangzhou.aliyuncs.com","Code":"Order.Opend","Message":"You have already opend the Cloud Governance Center service. Go to the console to start using it."}`, func(action string, _ map[string]string) {
			actions = append(actions, action)
		})

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudGovernanceService().Schema, map[string]interface{}{
			"enable": "On",
		})
		if err := dataSourceAliCloudGovernanceServiceRead(d, testUnitGovernanceServiceMockMeta(t, server)); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"OpenGovernanceService"}, actions)
		assert.Equal(t, "GovernanceServiceHasBeenOpened", d.Id())
		assert.Equal(t, "Opened", d.Get("status"))
	})

	t.Run("ErrorPropagated", func(t *testing.T) {
		server := testUnitGovernanceServiceStubServer(t, http.StatusNotFound, `{"RequestId":"stub","HostId":"governance.cn-hangzhou.aliyuncs.com","Code":"InvalidEnterpriseRealName.NotFound","Message":"The enterprise real-name verification has not been completed."}`, func(_ string, _ map[string]string) {})

		d := schema.TestResourceDataRaw(t, dataSourceAliCloudGovernanceService().Schema, map[string]interface{}{
			"enable": "On",
		})
		err := dataSourceAliCloudGovernanceServiceRead(d, testUnitGovernanceServiceMockMeta(t, server))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "InvalidEnterpriseRealName.NotFound")
	})
}
