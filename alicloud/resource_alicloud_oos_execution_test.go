package alicloud

import (
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
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_oos_execution", &resource.Sweeper{
		Name: "alicloud_oos_execution",
		F:    testSweepOosExecution,
	})
}

func testSweepOosExecution(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	action := "ListExecutions"
	for {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "testSweepOosExecution", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Executions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Executions", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			name := item["TemplateName"]
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping OOS Execution: %s", item["ExecutionId"].(string))
				continue
			}
			action = "DeleteExecutions"
			request := map[string]interface{}{
				"ExecutionIds": convertListToJsonString(convertListStringToListInterface([]string{item["ExecutionId"].(string)})),
			}
			_, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete OOS Execution (%s): %s", item["ExecutionId"].(string), err)
			}
			// Waiting 5 seconds to ensure OOS Execution have been deleted.
			time.Sleep(5 * time.Second)
			log.Printf("[INFO] Delete OOS Execution success: %s ", item["ExecutionId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudOOSExecution_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_execution.default"
	ra := resourceAttrInit(resourceId, OosExecutionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosExecution%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OosExecutionBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"template_name": "${alicloud_oos_template.default.template_name}",
					"parameters":    `{\"Status\":\"Running\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name":  CHECKSET,
						"parameters":     CHECKSET,
						"mode":           "Automatic",
						"resource_ids.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudOOSExecution_complete(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_execution.default"
	ra := resourceAttrInit(resourceId, OosExecutionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosExecution%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OosExecutionBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"template_name":       "${alicloud_oos_template.default.template_name}",
					"description":         "From TF Test",
					"parameters":          `{\"Status\":\"Running\"}`,
					"mode":                "Automatic",
					"safety_check":        "Skip",
					"template_version":    "${alicloud_oos_template.default.template_version}",
					"loop_mode":           "",
					"parent_execution_id": "",
					"template_content":    "",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Execution Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name":       CHECKSET,
						"description":         "From TF Test",
						"parameters":          CHECKSET,
						"mode":                "Automatic",
						"safety_check":        "Skip",
						"template_version":    CHECKSET,
						"loop_mode":           "",
						"parent_execution_id": "",
						"template_content":    "",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Execution Test",
						// ListExecutions returns Counters as an object while the
						// schema is a string; Read must store the JSON encoding.
						"counters":       REGEXMATCH + "TotalTasks",
						"resource_ids.#": "1",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_ids.0", resourceId, "id"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudOOSExecution_resourceGroup(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_execution.default"
	ra := resourceAttrInit(resourceId, OosExecutionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosExecution%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OosExecutionResourceGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		// tags is Optional (not Computed) and only written by Read, so the
		// id-only refresh check sees tags.% = "0" that the applied state never
		// held. Ignore just that key so the new attributes stay checked.
		IDRefreshIgnore:   []string{"tags"},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"template_name":     "${alicloud_oos_template.default.template_name}",
					"parameters":        `{\"Status\":\"Running\"}`,
					"resource_group_id": "${alicloud_resource_manager_resource_group.default.id}",
					"safety_check":      "Skip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name":     CHECKSET,
						"parameters":        CHECKSET,
						"resource_group_id": CHECKSET,
						"mode":              "Automatic",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.default", "id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_name":     "${alicloud_oos_template.default.template_name}",
					"parameters":        `{\"Status\":\"Running\"}`,
					"resource_group_id": "${alicloud_resource_manager_resource_group.change.id}",
					"safety_check":      "Skip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "alicloud_resource_manager_resource_group.change", "id"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var OosExecutionMap = map[string]string{
	"create_date":       CHECKSET,
	"executed_by":       CHECKSET,
	"is_parent":         CHECKSET,
	"resource_group_id": CHECKSET,
	"start_date":        CHECKSET,
	"status":            CHECKSET,
	"template_id":       CHECKSET,
	"update_date":       CHECKSET,
}

func OosExecutionBasicdependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Describe instances of given status",
			"Parameters":{
			  "Status":{
				"Type": "String",
				"Description": "(Required) The status of the Ecs instance."
			  }
			},
			"Tasks": [
			  {
				"Properties" :{
				  "Parameters":{
					"Status": "{{ Status }}"
				  },
				  "API": "DescribeInstances",
				  "Service": "Ecs"
				},
				"Name": "foo",
				"Action": "ACS::ExecuteApi"
			  }]
		  }
		  EOF
		  template_name = "%s"
		  version_name = "test"
		  tags = {
			"Created" = "TF",
			"For" = "template Test"
		  }
		}
	`, name)
}

func OosExecutionResourceGroupDependence(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_resource_manager_resource_group" "default" {
  resource_group_name = "%s"
  display_name        = "%s"
}

resource "alicloud_resource_manager_resource_group" "change" {
  resource_group_name = "%s-c"
  display_name        = "%s-c"
}
`, OosExecutionBasicdependence(name), name, name, name, name)
}

// lintignore: R001
func TestUnitAlicloudOOSExecution(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_oos_execution"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_oos_execution"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"description":         "From TF Test",
		"parameters":          `{\"Status\":\"Running\"}`,
		"loop_mode":           "loop_mode",
		"mode":                "mode",
		"parent_execution_id": "parent_execution_id",
		"resource_group_id":   "resource_group_id",
		"safety_check":        "safety_check",
		"template_content":    "template_content",
		"template_name":       "template_name",
		"template_version":    "template_version",
		"tags": map[string]interface{}{
			"Created": "TF",
			"For":     "Test",
		},
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Executions": []interface{}{
			map[string]interface{}{
				"Counters":   map[string]interface{}{"SuccessTasks": 1, "FailedTasks": 0, "TotalTasks": 1},
				"CreateDate": "create_date",
				"EndDate":    "end_date",
				"ExecutedBy": "executed_by",
				"IsParent":   "is_parent",
				"Mode":       "mode",
				"Outputs":    "outputs",
				"Parameters": map[string]interface{}{"Status": "Running"},
				"Tags": map[string]interface{}{
					"Created": "TF",
					"For":     "Test",
				},
				"ParentExecutionId": "parent_execution_id",
				"ResourceGroupId":   "resource_group_id",
				"RamRole":           "ram_role",
				"StartDate":         "start_date",
				"Status":            "Success",
				"StatusMessage":     "status_message",
				"TemplateId":        "template_id",
				"TemplateName":      "template_name",
				"TemplateVersion":   "template_version",
				"UpdateDate":        "update_date",
				"ExecutionId":       "MockExecutionId",
			},
		},
		"Execution": map[string]interface{}{
			"ExecutionId": "MockExecutionId",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_oos_execution", "MockExecutionId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudOosExecutionCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudOosExecutionCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudOosExecutionCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Delete Method
	d.SetId("MockExecutionId")

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudOosExecutionDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudOosExecutionDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudOosExecutionDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeOosExecutionNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudOosExecutionRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadNormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		d.SetId("MockExecutionId")
		err := resourceAlicloudOosExecutionRead(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		// Counters is returned as an object by ListExecutions while the schema is
		// a string, so Read must write the JSON encoding of it.
		counters, ok := d.Get("counters").(string)
		assert.True(t, ok)
		assert.Contains(t, counters, "SuccessTasks")
		// resource_ids is derived from the execution tags read in the same call.
		assert.Equal(t, []interface{}{"MockExecutionId"}, d.Get("resource_ids").([]interface{}))
	})

	t.Run("ReadDescribeOosExecutionAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudOosExecutionRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// TestUnitAlicloudOosExecutionStartRequestWirePayload asserts the exact
// StartExecution payload the provider hands to the transport for a configured
// execution. It needs no cloud credentials, so the request-level truth for
// resource_group_id is reproducible in the changed-unit-test gate.
func TestUnitAlicloudOosExecutionStartRequestWirePayload(t *testing.T) {
	d, _ := schema.InternalMap(Provider().ResourcesMap["alicloud_oos_execution"].Schema).Data(nil, nil)
	assert.Nil(t, d.Set("template_name", "tpl-wire"))
	assert.Nil(t, d.Set("resource_group_id", "rg-wire"))
	assert.Nil(t, d.Set("parameters", `{"Status":"Running"}`))
	assert.Nil(t, d.Set("tags", map[string]interface{}{
		"k1": "v1",
	}))

	request, err := oosExecutionStartExecutionRequest(d, "cn-beijing")
	assert.Nil(t, err)
	assert.Equal(t, "cn-beijing", request["RegionId"])
	assert.Equal(t, "tpl-wire", request["TemplateName"])
	assert.Equal(t, "rg-wire", request["ResourceGroupId"])
	assert.Equal(t, `{"k1":"v1"}`, request["Tags"])
}

// TestUnitAlicloudOosExecutionReadMapsExecutionResponse asserts the mapping from
// the ListExecutions response object into state, including the derived
// resource_ids, without cloud credentials.
func TestUnitAlicloudOosExecutionReadMapsExecutionResponse(t *testing.T) {
	d, _ := schema.InternalMap(Provider().ResourcesMap["alicloud_oos_execution"].Schema).Data(nil, nil)
	d.SetId("exec-wire")

	// The execution object is what DescribeOosExecution returns for
	// ListExecutions with ExecutionId=exec-wire and a bound tag.
	setOosExecutionState(d, map[string]interface{}{
		"ExecutionId":     "exec-wire",
		"Status":          "Success",
		"ResourceGroupId": "rg-wire",
		"Tags": map[string]interface{}{
			"k1": "v1",
		},
	})

	assert.Equal(t, "rg-wire", d.Get("resource_group_id"))
	assert.Equal(t, []interface{}{"exec-wire"}, d.Get("resource_ids").([]interface{}))
	assert.Equal(t, map[string]interface{}{"k1": "v1"}, d.Get("tags"))

	// Without tags the derived resource_ids stays empty.
	untagged, _ := schema.InternalMap(Provider().ResourcesMap["alicloud_oos_execution"].Schema).Data(nil, nil)
	untagged.SetId("exec-wire-untagged")
	setOosExecutionState(untagged, map[string]interface{}{
		"ExecutionId": "exec-wire-untagged",
		"Status":      "Success",
		"Tags":        map[string]interface{}{},
	})
	assert.Equal(t, []interface{}{}, untagged.Get("resource_ids").([]interface{}))
}

// oosExecutionTestClient builds a real connectivity client whose OOS endpoint is
// a local HTTP server, so resource entry points can run without cloud access and
// the serialized request can be asserted. It follows the ecsSnapshotTestClient
// recipe used elsewhere in this package.
func oosExecutionTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
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
	endpoints.Store("oos", endpoint)
	return client
}

// TestUnitAlicloudOosExecutionEntryWire drives the public Create and Read entry
// points against a local HTTP server and asserts the serialized requests: the
// StartExecution payload carries the configured ResourceGroupId and Tags, and
// the ListExecutions response object is what populates resource_group_id and
// the derived resource_ids. It needs no cloud credentials.
func TestUnitAlicloudOosExecutionEntryWire(t *testing.T) {
	var mu sync.Mutex
	var startExecutionForm, listExecutionsForm url.Values

	client := oosExecutionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
			return
		}
		action := r.Form.Get("Action")
		mu.Lock()
		switch action {
		case "StartExecution":
			startExecutionForm = r.Form
		case "ListExecutions":
			listExecutionsForm = r.Form
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "StartExecution":
			fmt.Fprint(w, `{"Execution":{"ExecutionId":"exec-wire"}}`)
		case "ListExecutions":
			fmt.Fprint(w, `{"Executions":[{"ExecutionId":"exec-wire","Status":"Success","ResourceGroupId":"rg-wire","Tags":{"k1":"v1"}}]}`)
		case "DeleteExecutions":
			fmt.Fprint(w, `{}`)
		default:
			t.Errorf("unexpected OOS action %q", action)
			fmt.Fprint(w, `{}`)
		}
	})

	d, _ := schema.InternalMap(Provider().ResourcesMap["alicloud_oos_execution"].Schema).Data(nil, nil)
	assert.Nil(t, d.Set("template_name", "tpl-wire"))
	assert.Nil(t, d.Set("resource_group_id", "rg-wire"))
	assert.Nil(t, d.Set("parameters", `{"Status":"Running"}`))
	assert.Nil(t, d.Set("tags", map[string]interface{}{
		"k1": "v1",
	}))

	assert.Nil(t, resourceAlicloudOosExecutionCreate(d, client))
	assert.Equal(t, "exec-wire", d.Id())

	mu.Lock()
	startForm := startExecutionForm
	listForm := listExecutionsForm
	mu.Unlock()

	assert.Equal(t, "StartExecution", startForm.Get("Action"))
	assert.Equal(t, "rg-wire", startForm.Get("ResourceGroupId"))
	assert.Equal(t, "tpl-wire", startForm.Get("TemplateName"))
	assert.Equal(t, "cn-hangzhou", startForm.Get("RegionId"))
	assert.JSONEq(t, `{"k1":"v1"}`, startForm.Get("Tags"))

	assert.Equal(t, "ListExecutions", listForm.Get("Action"))
	assert.Equal(t, "exec-wire", listForm.Get("ExecutionId"))

	assert.Equal(t, "rg-wire", d.Get("resource_group_id"))
	assert.Equal(t, []interface{}{"exec-wire"}, d.Get("resource_ids").([]interface{}))
}
