package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_auto_snapshot_policy", &resource.Sweeper{
		Name: "alicloud_ecs_auto_snapshot_policy",
		F:    testSweepEcsAutoSnapshotPolicy,
	})
}

func testSweepEcsAutoSnapshotPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId
	for {
		action := "DescribeAutoSnapshotPolicyEx"
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AutoSnapshotPolicies.AutoSnapshotPolicy", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AutoSnapshotPolicyName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs SnapShot Policy: %s (%s)", item["AutoSnapshotPolicyName"], item["AutoSnapshotPolicyId"])
				continue
			}
			action = "DeleteAutoSnapshotPolicy"
			request := map[string]interface{}{
				"autoSnapshotPolicyId": item["AutoSnapshotPolicyId"],
				"regionId":             client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs SnapShot Policy (%s (%s)): %s", item["AutoSnapshotPolicyName"], item["InstanceId"], err)
				continue
			}
			log.Printf("[INFO] Delete Ecs SnapShot Policy success: %s ", item["AutoSnapshotPolicyId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestUnitAliCloudECSAutoSnapshotPolicy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ecs_auto_snapshot_policy"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ecs_auto_snapshot_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"name":                            "name",
		"repeat_weekdays":                 []string{"1"},
		"retention_days":                  -1,
		"time_points":                     []string{"1"},
		"copied_snapshots_retention_days": 2,
		"enable_cross_region_copy":        true,
		"target_copy_regions":             []string{"cn-beijing"},
		"tags": map[string]string{
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
		"AutoSnapshotPolicies": map[string]interface{}{
			"AutoSnapshotPolicy": []interface{}{
				map[string]interface{}{
					"CopiedSnapshotsRetentionDays": 2,
					"EnableCrossRegionCopy":        true,
					"AutoSnapshotPolicyName":       "name",
					"RepeatWeekdays":               "[\"1\"]",
					"RetentionDays":                -1,
					"Status":                       "Normal",
					"Tags": map[string]interface{}{
						"Created": "TF",
						"For":     "Test",
					},
					"TargetCopyRegions":    "[\"cn-beijing\"]",
					"TimePoints":           "[\"1\"]",
					"AutoSnapshotPolicyId": "MockAutoSnapshotPolicyId",
				},
			},
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_auto_snapshot_policy", "MockAutoSnapshotPolicyId"))
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
			result["AutoSnapshotPolicyId"] = "MockAutoSnapshotPolicyId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudEcsAutoSnapshotPolicyCreate(d, rawClient)
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
		err := resourceAliCloudEcsAutoSnapshotPolicyCreate(d, rawClient)
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
		err := resourceAliCloudEcsAutoSnapshotPolicyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateNonRetryableError", func(t *testing.T) {
		retryFlag := false
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
		err := resourceAliCloudEcsAutoSnapshotPolicyCreate(dCreate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("MockAutoSnapshotPolicyId"))

	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudEcsAutoSnapshotPolicyUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyAutoSnapshotPolicyExAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"copied_snapshots_retention_days", "enable_cross_region_copy", "name", "repeat_weekdays", "retention_days", "target_copy_regions", "time_points"} {
			switch p["alicloud_ecs_auto_snapshot_policy"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			case schema.TypeSet:
				diff.SetAttribute("repeat_weekdays.#", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("repeat_weekdays.0", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("target_copy_regions.#", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("target_copy_regions.0", &terraform.ResourceAttrDiff{Old: "", New: "cn-beijing"})
				diff.SetAttribute("time_points.#", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("time_points.0", &terraform.ResourceAttrDiff{Old: "", New: "1"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_auto_snapshot_policy"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudEcsAutoSnapshotPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyAutoSnapshotPolicyExNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"copied_snapshots_retention_days", "enable_cross_region_copy", "name", "repeat_weekdays", "retention_days", "target_copy_regions", "time_points", "tags"} {
			switch p["alicloud_ecs_auto_snapshot_policy"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			case schema.TypeSet:
				diff.SetAttribute("repeat_weekdays.#", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("repeat_weekdays.0", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("target_copy_regions.#", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("target_copy_regions.0", &terraform.ResourceAttrDiff{Old: "", New: "cn-beijing"})
				diff.SetAttribute("time_points.#", &terraform.ResourceAttrDiff{Old: "", New: "1"})
				diff.SetAttribute("time_points.0", &terraform.ResourceAttrDiff{Old: "", New: "1"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_auto_snapshot_policy"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudEcsAutoSnapshotPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudEcsAutoSnapshotPolicyDelete(d, rawClient)
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
		err := resourceAliCloudEcsAutoSnapshotPolicyDelete(d, rawClient)
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
		err := resourceAliCloudEcsAutoSnapshotPolicyDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeEcsAutoSnapshotPolicyNotFound", func(t *testing.T) {
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
		err := resourceAliCloudEcsAutoSnapshotPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeEcsAutoSnapshotPolicyAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudEcsAutoSnapshotPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Ecs AutoSnapshotPolicy. >>> Resource test cases, automatically generated.
// Case 自动快照策略 7893
func TestAccAliCloudECSAutoSnapshotPolicy_basic7893(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"1", "2"},
					"retention_days":  "1",
					"time_points":     []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "2",
						"retention_days":    "1",
						"time_points.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_snapshot_policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copied_snapshots_retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copied_snapshots_retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_region_copy": "true",
					"target_copy_regions":      []string{"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_region_copy": "true",
						"target_copy_regions.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copy_encryption_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"3", "4", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"7", "6", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudECSAutoSnapshotPolicy_basic7893_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays":                 []string{"1", "2"},
					"retention_days":                  "1",
					"time_points":                     []string{"1", "2"},
					"auto_snapshot_policy_name":       name,
					"copied_snapshots_retention_days": "2",
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
					"enable_cross_region_copy": "true",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"target_copy_regions":      []string{"cn-hangzhou"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#":               "2",
						"retention_days":                  "1",
						"time_points.#":                   "2",
						"auto_snapshot_policy_name":       name,
						"copied_snapshots_retention_days": "2",
						"copy_encryption_configuration.#": "1",
						"enable_cross_region_copy":        "true",
						"resource_group_id":               CHECKSET,
						"target_copy_regions.#":           "1",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudEcsAutoSnapshotPolicyMap7893 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudEcsAutoSnapshotPolicyBasicDependence7893(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = 7
  		key_state              = "Enabled"
	}
`, name)
}

// Case 适配废弃字段name 8035
func TestAccAliCloudECSAutoSnapshotPolicy_basic8035(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"1", "2"},
					"retention_days":  "1",
					"time_points":     []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "2",
						"retention_days":    "1",
						"time_points.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copied_snapshots_retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copied_snapshots_retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_region_copy": "true",
					"target_copy_regions":      []string{"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_region_copy": "true",
						"target_copy_regions.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copy_encryption_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"3", "4", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"7", "6", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudECSAutoSnapshotPolicy_basic8035_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays":                 []string{"1", "2"},
					"retention_days":                  "1",
					"time_points":                     []string{"1", "2"},
					"name":                            name,
					"copied_snapshots_retention_days": "2",
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
					"enable_cross_region_copy": "true",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"target_copy_regions":      []string{"cn-hangzhou"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#":               "2",
						"retention_days":                  "1",
						"time_points.#":                   "2",
						"name":                            name,
						"copied_snapshots_retention_days": "2",
						"copy_encryption_configuration.#": "1",
						"enable_cross_region_copy":        "true",
						"resource_group_id":               CHECKSET,
						"target_copy_regions.#":           "1",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test Ecs AutoSnapshotPolicy. <<< Resource test cases, automatically generated.
