package oos_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oos"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_oos_patch_baseline",
		&resource.Sweeper{
			Name: "alicloud_oos_patch_baseline",
			F:    testSweepOosPatchBaseline,
		})
}

func testSweepOosPatchBaseline(region string) error {
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.OOSSupportRegions) {
		log.Printf("[INFO] Skipping Oos Patch Baseline unsupported region: %s", region)
		return nil
	}

	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListPatchBaselines"
	request := map[string]interface{}{}
	request["MaxResults"] = helper.PageSizeLarge

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.PatchBaselines", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.PatchBaselines", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Oos Patch Baseline: %s", item["Name"].(string))
				continue
			}
			action := "DeletePatchBaseline"
			request := map[string]interface{}{
				"Name": item["Name"],
			}
			_, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Oos Patch Baseline (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Oos Patch Baseline success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudOOSPatchBaseline_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOOSPatchBaselineMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soospatchbaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOOSPatchBaselineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operation_system":    "Windows",
					"patch_baseline_name": "${var.name}",
					"approval_rules":      `{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operation_system":    "Windows",
						"patch_baseline_name": name,
						"approval_rules":      acctest.CHECKSET,
						"resource_group_id":   acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"description":       "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"approved_patches_enable_non_security": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approved_patches_enable_non_security": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"approved_patches_enable_non_security": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approved_patches_enable_non_security": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"approved_patches_enable_non_security": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approved_patches_enable_non_security": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"approved_patches": []string{
						"KB1", "KB5", "KB6"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approved_patches.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []string{
						"[main]", "name=MyCustomRepository", "name=MyCustomRepository", "enabled=1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"approval_rules": `{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Values\":[\"OS\"],\"Key\":\"PatchSet\"},{\"Values\":[[\"Windows\"]],\"Key\":\"ProductFamily\"},{\"Values\":[\"Windows 10\",\"Windows 7\"],\"Key\":\"Product\"},{\"Values\":[\"*\"],\"Key\":\"Classification\"},{\"Values\":[\"*\"],\"Key\":\"Severity\"}],\"ApproveAfterDays\":7,\"ComplianceLevel\":\"Unspecified\"},{\"PatchFilterGroup\":[{\"Values\":[\"Application\"],\"Key\":\"PatchSet\"},{\"Values\":[\"*\"],\"Key\":\"ProductFamily\"},{\"Values\":[\"*\"],\"Key\":\"Product\"},{\"Values\":[\"*\"],\"Key\":\"Classification\"},{\"Values\":[\"*\"],\"Key\":\"Severity\"}],\"ApproveAfterDays\":7,\"ComplianceLevel\":\"Unspecified\"}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approval_rules": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "${var.name}_update",
					"approval_rules": `{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name + "_update",
						"approval_rules": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"approved_patches_enable_non_security"},
			},
		},
	})
}

func TestAccAliCloudOOSPatchBaseline_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOOSPatchBaselineMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soospatchbaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOOSPatchBaselineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"operation_system":    "Windows",
					"patch_baseline_name": "${var.name}",
					"description":         "${var.name}",
					"approval_rules":      `{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operation_system":    "Windows",
						"patch_baseline_name": name,
						"description":         name,
						"approval_rules":      acctest.CHECKSET,
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

var AlicloudOOSPatchBaselineMap0 = map[string]string{
	"patch_baseline_name": acctest.CHECKSET,
}

func AlicloudOOSPatchBaselineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}

// lintignore: R001
func TestUnitAlicloudOOSPatchBaseline(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_oos_patch_baseline"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_oos_patch_baseline"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"operation_system":    "Windows",
		"patch_baseline_name": "MockName",
		"description":         "description",
		"approval_rules":      `{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}`,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"PatchBaseline": map[string]interface{}{
			"Name":            "MockName",
			"Description":     "description",
			"OperationSystem": "Windows",
			"ApprovalRules": map[string]interface{}{
				"PatchRules": []interface{}{
					map[string]interface{}{
						"EnableNonSecurity": true,
					},
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_oos_patch_baseline", "MockName"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := oos.ResourceAliCloudOosPatchBaselineCreate(d, rawClient)
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
		err := oos.ResourceAliCloudOosPatchBaselineCreate(d, rawClient)
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
		err := oos.ResourceAliCloudOosPatchBaselineCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockName")

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := oos.ResourceAliCloudOosPatchBaselineUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateParameterAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"approval_rules", "description"} {
			switch p["alicloud_oos_patch_baseline"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_oos_patch_baseline"].Schema).Data(nil, diff)
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
		err := oos.ResourceAliCloudOosPatchBaselineUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateParameterNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"approval_rules", "description"} {
			switch p["alicloud_oos_patch_baseline"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_oos_patch_baseline"].Schema).Data(nil, diff)
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
		err := oos.ResourceAliCloudOosPatchBaselineUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := oos.ResourceAliCloudOosPatchBaselineDelete(d, rawClient)
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
		err := oos.ResourceAliCloudOosPatchBaselineDelete(d, rawClient)
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
		err := oos.ResourceAliCloudOosPatchBaselineDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeOosPatchBaselineNotFound", func(t *testing.T) {
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
		err := oos.ResourceAliCloudOosPatchBaselineRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeOosPatchBaselineAbnormal", func(t *testing.T) {
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
		err := oos.ResourceAliCloudOosPatchBaselineRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Case OOS patch baseline接入Tag_ubuntu 5946
func TestAccAliCloudOosPatchBaseline_basic5946(t *testing.T) {
	t.Skipf("There is currently a delay in the tag service, causing this use case to fail. However, the functionality is actually working normally.")
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap5946)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence5946)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "Ubuntu",
					"rejected_patches": []string{
						"KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"approved_patches": []string{
						"KB6", "MS7", "MS5"},
					"approved_patches_enable_non_security": "true",
					"sources": []string{
						"new"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "ALLOW_AS_DEPENDENCY",
						"description":                          "xxxx",
						"patch_baseline_name":                  name,
						"approval_rules":                       acctest.CHECKSET,
						"operation_system":                     "Ubuntu",
						"rejected_patches.#":                   "2",
						"resource_group_id":                    acctest.CHECKSET,
						"approved_patches.#":                   "3",
						"approved_patches_enable_non_security": "true",
						"sources.#":                            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"approved_patches": []string{
						"MS7"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
						"resource_group_id":       acctest.CHECKSET,
						"approved_patches.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
					"approved_patches": []string{},
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux4Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux5Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux6Source\\\"}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
						"approved_patches.#": "0",
						"sources.#":          "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                          "123213",
					"rejected_patches":                     []string{},
					"approved_patches_enable_non_security": "false",
					"sources":                              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                          acctest.CHECKSET,
						"rejected_patches.#":                   "0",
						"approved_patches_enable_non_security": "false",
						"sources.#":                            "0",
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudOosPatchBaselineMap5946 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence5946(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case OOS patch baseline接入Tag_aliyun_linux 5945
func TestAccAliCloudOosPatchBaseline_basic5945(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap5945)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence5945)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "AliyunLinux",
					"rejected_patches": []string{
						"KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"approved_patches": []string{
						"KB6", "MS7", "MS5"},
					"approved_patches_enable_non_security": "true",
					"sources": []string{
						"new"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "ALLOW_AS_DEPENDENCY",
						"description":                          "xxxx",
						"patch_baseline_name":                  name,
						"approval_rules":                       acctest.CHECKSET,
						"operation_system":                     "AliyunLinux",
						"rejected_patches.#":                   "2",
						"resource_group_id":                    acctest.CHECKSET,
						"approved_patches.#":                   "3",
						"approved_patches_enable_non_security": "true",
						"sources.#":                            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"approved_patches": []string{
						"MS7"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
						"resource_group_id":       acctest.CHECKSET,
						"approved_patches.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
					"approved_patches": []string{},
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux4Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux5Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux6Source\\\"}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
						"approved_patches.#": "0",
						"sources.#":          "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                          "123213",
					"rejected_patches":                     []string{},
					"approved_patches_enable_non_security": "false",
					"sources":                              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                          acctest.CHECKSET,
						"rejected_patches.#":                   "0",
						"approved_patches_enable_non_security": "false",
						"sources.#":                            "0",
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

var AlicloudOosPatchBaselineMap5945 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence5945(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case OOS patch baseline接入Tag_debian 5947
func TestAccAliCloudOosPatchBaseline_basic5947(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap5947)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence5947)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "Debian",
					"rejected_patches": []string{
						"KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"approved_patches": []string{
						"KB6", "MS7", "MS5"},
					"approved_patches_enable_non_security": "true",
					"sources": []string{
						"new"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "ALLOW_AS_DEPENDENCY",
						"description":                          "xxxx",
						"patch_baseline_name":                  name,
						"approval_rules":                       acctest.CHECKSET,
						"operation_system":                     "Debian",
						"rejected_patches.#":                   "2",
						"resource_group_id":                    acctest.CHECKSET,
						"approved_patches.#":                   "3",
						"approved_patches_enable_non_security": "true",
						"sources.#":                            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"approved_patches": []string{
						"MS7"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
						"resource_group_id":       acctest.CHECKSET,
						"approved_patches.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
					"approved_patches": []string{},
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux4Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux5Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux6Source\\\"}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
						"approved_patches.#": "0",
						"sources.#":          "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                          "123213",
					"rejected_patches":                     []string{},
					"approved_patches_enable_non_security": "false",
					"sources":                              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                          acctest.CHECKSET,
						"rejected_patches.#":                   "0",
						"approved_patches_enable_non_security": "false",
						"sources.#":                            "0",
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

var AlicloudOosPatchBaselineMap5947 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence5947(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case OOS patch baseline接入Tag_centos 5944
func TestAccAliCloudOosPatchBaseline_basic5944(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap5944)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence5944)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "CentOS",
					"rejected_patches": []string{
						"KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"approved_patches": []string{
						"KB6", "MS7", "MS5"},
					"approved_patches_enable_non_security": "true",
					"sources": []string{
						"new"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "ALLOW_AS_DEPENDENCY",
						"description":                          "xxxx",
						"patch_baseline_name":                  name,
						"approval_rules":                       acctest.CHECKSET,
						"operation_system":                     "CentOS",
						"rejected_patches.#":                   "2",
						"resource_group_id":                    acctest.CHECKSET,
						"approved_patches.#":                   "3",
						"approved_patches_enable_non_security": "true",
						"sources.#":                            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"Product\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"approved_patches": []string{
						"MS7"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
						"resource_group_id":       acctest.CHECKSET,
						"approved_patches.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
					"approved_patches": []string{},
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux4Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux5Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux6Source\\\"}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
						"approved_patches.#": "0",
						"sources.#":          "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                          "123213",
					"rejected_patches":                     []string{},
					"approved_patches_enable_non_security": "false",
					"sources":                              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                          acctest.CHECKSET,
						"rejected_patches.#":                   "0",
						"approved_patches_enable_non_security": "false",
						"sources.#":                            "0",
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

var AlicloudOosPatchBaselineMap5944 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence5944(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case OOS patch baseline接入Tag_windows 4719
func TestAccAliCloudOosPatchBaseline_basic4719(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap4719)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoos%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence4719)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "Windows",
					"rejected_patches": []string{
						"KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"approved_patches": []string{
						"KB6", "MS7", "MS5"},
					"approved_patches_enable_non_security": "true",
					"sources": []string{
						"new"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "ALLOW_AS_DEPENDENCY",
						"description":                          "xxxx",
						"patch_baseline_name":                  name,
						"approval_rules":                       acctest.CHECKSET,
						"operation_system":                     "Windows",
						"rejected_patches.#":                   "2",
						"resource_group_id":                    acctest.CHECKSET,
						"approved_patches.#":                   "3",
						"approved_patches_enable_non_security": "true",
						"sources.#":                            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"approved_patches": []string{
						"MS7"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
						"resource_group_id":       acctest.CHECKSET,
						"approved_patches.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
					"approved_patches": []string{},
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux4Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux5Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux6Source\\\"}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
						"approved_patches.#": "0",
						"sources.#":          "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                          "123213",
					"rejected_patches":                     []string{},
					"approved_patches_enable_non_security": "false",
					"sources":                              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                          acctest.CHECKSET,
						"rejected_patches.#":                   "0",
						"approved_patches_enable_non_security": "false",
						"sources.#":                            "0",
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

var AlicloudOosPatchBaselineMap4719 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence4719(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case OOS patch baseline接入Tag_new_3 3681
func TestAccAliCloudOosPatchBaseline_basic3681(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap3681)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoos%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence3681)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "Windows",
					"rejected_patches": []string{
						"KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux3Source\\\"}"},
					"approved_patches": []string{
						"KB1"},
					"approved_patches_enable_non_security": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "ALLOW_AS_DEPENDENCY",
						"description":                          "xxxx",
						"patch_baseline_name":                  name,
						"approval_rules":                       acctest.CHECKSET,
						"operation_system":                     "Windows",
						"rejected_patches.#":                   "2",
						"resource_group_id":                    acctest.CHECKSET,
						"sources.#":                            "1",
						"approved_patches.#":                   "1",
						"approved_patches_enable_non_security": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux2Source\\\"}"},
					"approved_patches_enable_non_security": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action":              "BLOCK",
						"description":                          "qwqww1",
						"approval_rules":                       acctest.CHECKSET,
						"rejected_patches.#":                   "4",
						"resource_group_id":                    acctest.CHECKSET,
						"sources.#":                            "1",
						"approved_patches_enable_non_security": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
					"sources": []string{
						"{\\\"Name\\\":\\\"AliyunLinux4Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux5Source\\\"}", "{\\\"Name\\\":\\\"AliyunLinux6Source\\\"}"},
					"approved_patches": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
						"sources.#":          "3",
						"approved_patches.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "123213",
					"rejected_patches": []string{},
					"sources":          []string{},
					"approved_patches": []string{
						"KB1", "KB5", "KB6"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        acctest.CHECKSET,
						"rejected_patches.#": "0",
						"sources.#":          "0",
						"approved_patches.#": "3",
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

var AlicloudOosPatchBaselineMap3681 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence3681(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case OOS patch baseline接入Tag_new 3648
func TestAccAliCloudOosPatchBaseline_basic3648(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap3648)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoos%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence3648)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "Windows",
					"rejected_patches": []string{
						"KB3", "KB4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
						"description":             "xxxx",
						"patch_baseline_name":     name,
						"approval_rules":          acctest.CHECKSET,
						"operation_system":        "Windows",
						"rejected_patches.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{
						"KB1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "123213",
					"rejected_patches": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        acctest.CHECKSET,
						"rejected_patches.#": "0",
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

var AlicloudOosPatchBaselineMap3648 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence3648(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case OOS patch baseline接入Tag 3558
func TestAccAliCloudOosPatchBaseline_basic3558(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_patch_baseline.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOosPatchBaselineMap3558)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oos.NewOosServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOosPatchBaseline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoos%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOosPatchBaselineBasicDependence3558)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
					"description":             "xxxx",
					"patch_baseline_name":     name,
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":true,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"operation_system":        "Windows",
					"rejected_patches": []string{
						"KB3", "KB4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "ALLOW_AS_DEPENDENCY",
						"description":             "xxxx",
						"patch_baseline_name":     name,
						"approval_rules":          acctest.CHECKSET,
						"operation_system":        "Windows",
						"rejected_patches.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches_action": "BLOCK",
					"description":             "qwqww1",
					"approval_rules":          "{\\\"PatchRules\\\":[{\\\"EnableNonSecurity\\\":false,\\\"PatchFilterGroup\\\":[{\\\"Values\\\":[\\\"OS\\\"],\\\"Key\\\":\\\"PatchSet\\\"}],\\\"ApproveAfterDays\\\":7,\\\"ComplianceLevel\\\":\\\"Medium\\\"}]}",
					"rejected_patches": []string{
						"KB1", "KB2", "KB3", "KB4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches_action": "BLOCK",
						"description":             "qwqww1",
						"approval_rules":          acctest.CHECKSET,
						"rejected_patches.#":      "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rejected_patches": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rejected_patches.#": "0",
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

var AlicloudOosPatchBaselineMap3558 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOosPatchBaselineBasicDependence3558(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
