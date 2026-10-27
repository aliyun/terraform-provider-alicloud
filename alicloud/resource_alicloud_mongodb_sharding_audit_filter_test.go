package alicloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
)

// lintignore: AT001
func TestAccAliCloudMongodbShardingAuditFilter_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_audit_filter.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudMongoDBShardingAuditFilterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbshardingauditfilter-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBShardingAuditFilterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_mongodb_sharding_instance.default.id}",
					"audit_status":            "enable",
					"filter":                  "admin,slow",
					"role_type":               "mongos",
					"service_type":            "V2_Standard",
					"storage_period":          "30",
					"hot_storage_period":      "7",
					"audit_log_switch_source": "Console",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"audit_status":            "enable",
						"filter":                  "admin,slow",
						"role_type":               "mongos",
						"service_type":            "V2_Standard",
						"storage_period":          "30",
						"hot_storage_period":      "7",
						"audit_log_switch_source": "Console",
					}),
				),
			},
			{
				// A real policy change while audit_log_switch_source stays explicitly configured: the
				// request carries the configured source together with the policy change.
				Config: testAccConfig(map[string]interface{}{
					"audit_status":       "disabled",
					"filter":             "update,insert,delete",
					"role_type":          "db",
					"service_type":       "V2_Standard",
					"storage_period":     "180",
					"hot_storage_period": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status":            "disabled",
						"filter":                  "update,insert,delete",
						"role_type":               "db",
						"storage_period":          "180",
						"hot_storage_period":      "5",
						"audit_log_switch_source": "Console",
					}),
				),
			},
			{
				// Turning the audit log back on is a real policy change that omits the optional
				// audit_log_switch_source: the request must be built without the parameter and the
				// attribute must leave the state instead of being held at the earlier explicit value.
				Config: testAccConfig(map[string]interface{}{
					"audit_status":            "enable",
					"filter":                  "admin,slow",
					"role_type":               "mongos",
					"service_type":            "V2_Standard",
					"audit_log_switch_source": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status":            "enable",
						"filter":                  "admin,slow",
						"role_type":               "mongos",
						"audit_log_switch_source": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// The write-only audit_log_switch_source is never returned by the query APIs, and
				// service_type/role_type are Optional+Computed values the read APIs cannot restore as
				// declared.
				ImportStateVerifyIgnore: []string{"service_type", "role_type", "audit_log_switch_source"},
			},
			{
				// Import verification lives in the step above: ImportStateVerify compares the imported
				// state with the pre-import state, ignoring only the write-only
				// audit_log_switch_source and the Optional+Computed service_type/role_type, which the
				// read APIs cannot restore as declared. That import state is discarded after the step,
				// so this replay runs against the state of the re-enable step and pins that turning the
				// audit log back on converges without drift.
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "enable",
					"filter":       "admin,slow",
					"role_type":    "mongos",
					"service_type": "V2_Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "enable",
						"filter":       "admin,slow",
						"role_type":    "mongos",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
			},
			{
				// Changing a real attribute stays an in-place update (the only ForceNew field,
				// db_instance_id, is unchanged), and the framework then requires the follow-up plan to
				// be empty. The filter set and role type mirror step 2, a combination the API accepts
				// on this instance.
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "enable",
					"filter":       "update,insert,delete",
					"role_type":    "db",
					"service_type": "V2_Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "enable",
						"filter":       "update,insert,delete",
						"role_type":    "db",
					}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudMongodbShardingAuditFilter_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_audit_filter.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudMongoDBShardingAuditFilterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbshardingauditfilter-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBShardingAuditFilterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Create without audit_log_switch_source: the optional request source has no schema
				// default, so the audit policy request leaves the parameter out.
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":     "${alicloud_mongodb_sharding_instance.default.id}",
					"audit_status":       "enable",
					"filter":             "admin,slow",
					"role_type":          "mongos",
					"storage_period":     "30",
					"service_type":       "V2_Standard",
					"hot_storage_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":     CHECKSET,
						"audit_status":       "enable",
						"filter":             "admin,slow",
						"role_type":          "mongos",
						"storage_period":     "30",
						"service_type":       "V2_Standard",
						"hot_storage_period": "7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_type", "role_type", "audit_log_switch_source"},
			},
		},
	})
}

var AliCloudMongoDBShardingAuditFilterMap0 = map[string]string{
	"storage_period":     CHECKSET,
	"filter":             CHECKSET,
	"service_type":       "V2_Standard",
	"hot_storage_period": CHECKSET,
}

func AliCloudMongoDBShardingAuditFilterBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = data.alicloud_vswitches.default.ids[0]
  engine_version = "4.2"
  name           = var.name
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
}
`, name)
}

// lintignore: R001
func TestUnitAliCloudMongodbShardingAuditFilter(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mongodb_sharding_audit_filter"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mongodb_sharding_audit_filter"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_id":          "CreateValue",
		"audit_status":            "enable",
		"filter":                  "admin,slow",
		"audit_log_switch_source": "Console",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeAuditPolicy
		"DBInstanceId":   "CreateValue",
		"LogAuditStatus": "Enable",
		"Filter":         "admin,slow",
		"RoleType":       "primary",
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceStatus": "Running",
					"DBInstanceId":     "CreateValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_mongodb_sharding_audit_filter", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudMongodbShardingAuditFilterCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"DBInstanceId":   "CreateValue",
		"LogAuditStatus": "Enable",
		"Filter":         "admin,slow",
		"RoleType":       "primary",
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceStatus": "Running",
					"DBInstanceId":     "CreateValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyAuditPolicy" {
				assert.Equal(t, "Console", body["AuditLogSwitchSource"])
			}
			if *action == "ModifyAuditPolicy" || *action == "ModifyAuditLogFilter" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudMongodbShardingAuditFilterCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_sharding_audit_filter"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudMongodbShardingAuditFilterUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"audit_status":   "disabled",
		"filter":         "update,insert,delete",
		"storage_period": 30,
	}
	diff, err := newInstanceDiff("alicloud_mongodb_sharding_audit_filter", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_sharding_audit_filter"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"LogAuditStatus": "Disabled",
		"Filter":         "update,insert,delete",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyAuditPolicy" {
				assert.Equal(t, "Console", body["AuditLogSwitchSource"])
			}
			if *action == "ModifyAuditPolicy" || *action == "ModifyAuditLogFilter" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudMongodbShardingAuditFilterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_sharding_audit_filter"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeAuditPolicy" || *action == "DescribeMongoDBLogConfig" || *action == "DescribeAuditLogFilter" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudMongodbShardingAuditFilterRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAliCloudMongodbShardingAuditFilterDelete(dExisted, rawClient)
	assert.Nil(t, err)
}

// TestUnitAliCloudMongodbShardingAuditFilterAuditLogSwitchSource captures the serialized
// ModifyAuditPolicy request at the HTTP boundary, so the write-only AuditLogSwitchSource is checked
// without a credential-bound mock that skips inside CI containers. The attribute is optional request
// attribution: it must be sent when configured and left out of the request when it is not, on both
// Create and a policy-changing Update.
func TestUnitAliCloudMongodbShardingAuditFilterAuditLogSwitchSource(t *testing.T) {
	var mu sync.Mutex
	var policyForms []url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %v", err)
		}
		if r.Form.Get("Action") == "ModifyAuditPolicy" {
			mu.Lock()
			policyForms = append(policyForms, r.PostForm)
			mu.Unlock()
			// Non-retryable error: the request is already serialized on the wire, and the error
			// returns before the hard-coded five-minute state-change delay.
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"RequestId":"mock","Code":"InvalidParameter","Message":"mock request boundary"}`)
			return
		}
		fmt.Fprint(w, `{"RequestId":"mock"}`)
	}))
	t.Cleanup(server.Close)

	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("dds", endpoint)

	res := resourceAliCloudMongodbShardingAuditFilter()

	// Create with an explicit source.
	createExplicit := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":          "dds-mock",
		"audit_status":            "enable",
		"filter":                  "admin,slow",
		"audit_log_switch_source": "Console",
	})
	assert.Error(t, resourceAliCloudMongodbShardingAuditFilterCreate(createExplicit, client))

	// Create without the attribute: there is no schema default, so the request must leave the
	// parameter out instead of sending an empty or provider-chosen value.
	createOmitted := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id": "dds-mock",
		"audit_status":   "enable",
		"filter":         "admin,slow",
	})
	assert.Equal(t, "", createOmitted.Get("audit_log_switch_source"))
	assert.Error(t, resourceAliCloudMongodbShardingAuditFilterCreate(createOmitted, client))

	// Policy-changing Update with the attribute configured.
	stateData, _ := schema.InternalMap(res.Schema).Data(nil, nil)
	assert.NoError(t, stateData.Set("db_instance_id", "dds-mock"))
	assert.NoError(t, stateData.Set("audit_status", "enable"))
	assert.NoError(t, stateData.Set("filter", "admin,slow"))
	assert.NoError(t, stateData.Set("audit_log_switch_source", "Console"))
	stateData.SetId("dds-mock")
	baseState := stateData.State()
	attributes := map[string]interface{}{
		"db_instance_id":          "dds-mock",
		"audit_status":            "enable",
		"filter":                  "admin,slow",
		"audit_log_switch_source": "Console",
	}
	diff, err := newInstanceDiff("alicloud_mongodb_sharding_audit_filter", attributes, map[string]interface{}{"audit_status": "disabled"}, baseState)
	assert.NoError(t, err)
	updateData, _ := schema.InternalMap(res.Schema).Data(baseState, diff)
	assert.Error(t, resourceAliCloudMongodbShardingAuditFilterUpdate(updateData, client))

	// Policy-changing Update without the attribute: the state carries no source, so the request must
	// leave the parameter out even though the audit policy really changes.
	omittedStateData, _ := schema.InternalMap(res.Schema).Data(nil, nil)
	assert.NoError(t, omittedStateData.Set("db_instance_id", "dds-mock"))
	assert.NoError(t, omittedStateData.Set("audit_status", "enable"))
	assert.NoError(t, omittedStateData.Set("filter", "admin,slow"))
	omittedStateData.SetId("dds-mock")
	omittedState := omittedStateData.State()
	omittedAttributes := map[string]interface{}{
		"db_instance_id": "dds-mock",
		"audit_status":   "enable",
		"filter":         "admin,slow",
	}
	diff, err = newInstanceDiff("alicloud_mongodb_sharding_audit_filter", omittedAttributes, map[string]interface{}{"audit_status": "disabled"}, omittedState)
	assert.NoError(t, err)
	omittedUpdateData, _ := schema.InternalMap(res.Schema).Data(omittedState, diff)
	assert.Error(t, resourceAliCloudMongodbShardingAuditFilterUpdate(omittedUpdateData, client))

	mu.Lock()
	defer mu.Unlock()
	if !assert.Len(t, policyForms, 4) {
		return
	}
	assertSource := func(form url.Values, want, msg string) {
		values, ok := form["AuditLogSwitchSource"]
		if want == "" {
			assert.False(t, ok, "%s: the unconfigured parameter must be absent from the request", msg)
			return
		}
		assert.True(t, ok, "%s: the configured parameter must be sent", msg)
		assert.Equal(t, []string{want}, values, msg)
	}
	assertSource(policyForms[0], "Console", "create with the attribute configured")
	assertSource(policyForms[1], "", "create without the attribute")
	assertSource(policyForms[2], "Console", "policy-changing update with the attribute configured")
	assertSource(policyForms[3], "", "policy-changing update without the attribute")
	for i, form := range policyForms {
		assert.Equal(t, "dds-mock", form.Get("DBInstanceId"), "ModifyAuditPolicy request %d", i)
	}
}

// TestUnitAliCloudMongodbShardingAuditFilterStorageValues pins the edition-specific retention fields
// returned by DescribeMongoDBLogConfig: a V2_Standard audit log keeps the configured StoragePeriod in
// TtlForV2Standard, while TtlForStandard holds the untouched Standard-edition setting.
func TestUnitAliCloudMongodbShardingAuditFilterStorageValues(t *testing.T) {
	cases := []struct {
		name     string
		raw      map[string]interface{}
		reported string
		declared string
		storage  interface{}
		hot      interface{}
	}{
		{
			name:     "v2 standard reads cold and hot storage",
			raw:      map[string]interface{}{"TtlForStandard": 30, "TtlForV2Standard": 180, "HotTtlForV2Standard": 5},
			reported: "V2_Standard",
			declared: "V2_Standard",
			storage:  180,
			hot:      5,
		},
		{
			name:     "standard reads retention only",
			raw:      map[string]interface{}{"TtlForStandard": 180, "TtlForV2Standard": 30},
			reported: "Standard",
			declared: "Standard",
			storage:  180,
			hot:      nil,
		},
		{
			name:     "disabled reports trial so the declared edition decides",
			raw:      map[string]interface{}{"TtlForStandard": 30, "TtlForV2Standard": 180, "HotTtlForV2Standard": 5},
			reported: "Trial",
			declared: "V2_Standard",
			storage:  180,
			hot:      5,
		},
		{
			name:     "unknown edition without declaration falls back to standard",
			raw:      map[string]interface{}{"TtlForStandard": 30},
			reported: "Trial",
			declared: "",
			storage:  30,
			hot:      nil,
		},
		{
			name:     "omitted fields keep the state value",
			raw:      map[string]interface{}{},
			reported: "Trial",
			declared: "V2_Standard",
			storage:  nil,
			hot:      nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			storage, hot := mongodbShardingAuditFilterStorageValues(c.raw, c.reported, c.declared)
			assert.Equal(t, c.storage, storage)
			assert.Equal(t, c.hot, hot)
		})
	}
}

// TestUnitAliCloudMongodbShardingAuditFilterReadKeepsConfiguredSource pins that Read leaves an
// explicitly configured audit_log_switch_source untouched: none of the describe APIs return the
// write-only request parameter, so Read must neither clear nor overwrite it.
func TestUnitAliCloudMongodbShardingAuditFilterReadKeepsConfiguredSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %v", err)
		}
		switch r.Form.Get("Action") {
		case "DescribeMongoDBLogConfig":
			fmt.Fprint(w, `{"RequestId":"mock","ServiceType":"V2_Standard","TtlForStandard":30,"TtlForV2Standard":180,"HotTtlForV2Standard":5}`)
		case "DescribeAuditPolicy":
			fmt.Fprint(w, `{"RequestId":"mock","LogAuditStatus":"Enable"}`)
		case "DescribeAuditLogFilter":
			fmt.Fprint(w, `{"RequestId":"mock","Filter":"admin,slow","RoleType":"logic"}`)
		default:
			fmt.Fprint(w, `{"RequestId":"mock"}`)
		}
	}))
	t.Cleanup(server.Close)

	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("dds", endpoint)

	res := resourceAliCloudMongodbShardingAuditFilter()
	data := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_instance_id":          "dds-mock",
		"audit_status":            "enable",
		"filter":                  "admin,slow",
		"audit_log_switch_source": "Console",
	})
	data.SetId("dds-mock")

	assert.NoError(t, resourceAliCloudMongodbShardingAuditFilterRead(data, client))
	assert.Equal(t, "Console", data.Get("audit_log_switch_source"))
	assert.Equal(t, "enable", data.Get("audit_status"))
}

// Test Mongodb ShardingAuditFilter. >>> Resource test cases, hand-written.

// TestUnitAliCloudMongodbShardingAuditFilterDisabledImportPlan pins the imported-state transition of
// the optional write-only audit_log_switch_source. After importing an instance whose audit log is
// disabled the state has no service_type (Read keeps the declared value only while enabled) and the
// disabled edition suppression hides the declared V2_Standard. Because the attribute has no schema
// default, importing must not plan any change for it, and the plan must therefore carry no audit
// policy write that would fall back to ServiceType=Standard. The plan is built through the real SDK
// diff here, Update runs against it, and the follow-up plan must be empty.
func TestUnitAliCloudMongodbShardingAuditFilterDisabledImportPlan(t *testing.T) {
	var mu sync.Mutex
	var policyCalls, filterCalls int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyAuditPolicy":
			mu.Lock()
			policyCalls++
			mu.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"RequestId":"mock","Code":"InvalidParameter","Message":"unexpected policy write"}`)
		case "ModifyAuditLogFilter":
			mu.Lock()
			filterCalls++
			mu.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"RequestId":"mock","Code":"InvalidParameter","Message":"unexpected filter write"}`)
		case "DescribeMongoDBLogConfig":
			// Disabled audit logs report the Trial edition and no V2 retention fields.
			fmt.Fprint(w, `{"RequestId":"mock","ServiceType":"Trial","TtlForStandard":30}`)
		case "DescribeAuditPolicy":
			fmt.Fprint(w, `{"RequestId":"mock","LogAuditStatus":"Disabled"}`)
		case "DescribeAuditLogFilter":
			fmt.Fprint(w, `{"RequestId":"mock","Filter":"admin,slow","RoleType":"logic"}`)
		default:
			fmt.Fprint(w, `{"RequestId":"mock"}`)
		}
	}))
	t.Cleanup(server.Close)

	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("dds", endpoint)

	res := resourceAliCloudMongodbShardingAuditFilter()
	// State produced by terraform import: no service_type while the audit log is disabled and no
	// write-only audit_log_switch_source.
	prior := &terraform.InstanceState{ID: "dds-mock", Attributes: map[string]string{
		"id":             "dds-mock",
		"db_instance_id": "dds-mock",
		"audit_status":   "disabled",
		"filter":         "admin,slow",
	}}
	declared := map[string]interface{}{
		"db_instance_id": "dds-mock",
		"audit_status":   "disabled",
		"filter":         "admin,slow",
		"service_type":   "V2_Standard",
	}

	diff, err := res.Diff(context.Background(), prior, terraform.NewResourceConfigRaw(declared), nil)
	if !assert.NoError(t, err) {
		return
	}
	// Without a schema default for the write-only source the only difference against the declared
	// configuration is the suppressed disabled edition, so the plan must be empty: no attribute may
	// change and no audit policy write may be planned for the imported state.
	if diff != nil {
		for name, attribute := range diff.Attributes {
			if attribute.NewComputed {
				continue
			}
			assert.Failf(t, "unexpected planned change", "attribute %q would change", name)
		}
	}

	data, err := schema.InternalMap(res.Schema).Data(prior, diff)
	if !assert.NoError(t, err) {
		return
	}
	assert.NoError(t, resourceAliCloudMongodbShardingAuditFilterUpdate(data, client))

	mu.Lock()
	assert.Zero(t, policyCalls, "the source default alone must not rewrite the audit policy")
	assert.Zero(t, filterCalls, "no filter write is expected when only the source default is planned")
	mu.Unlock()

	followUp, err := res.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(declared), nil)
	assert.NoError(t, err)
	if followUp != nil {
		assert.Empty(t, followUp.Attributes, "the follow-up plan must converge")
	}
}
