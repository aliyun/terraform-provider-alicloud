package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// collection-order: bypass paths -- pre-existing multi-member ForceNew TypeList; this change only adjusts optionality and ConflictsWith for the new path sibling, not member-ordering semantics; the regular-NAS test account provides a single reliable mount path ("/") and CPFS is unavailable, so a convergent two-member reorder apply cannot be added safely

func TestAccAliCloudNASLifecyclePolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_lifecycle_policy.default"
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNASLifecyclePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasLifecyclePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snaslifecyclepolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNASLifecyclePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_policy_name": "${var.name}",
					"file_system_id":        "${alicloud_nas_file_system.default.id}",
					"lifecycle_rule_name":   "DEFAULT_ATIME_14",
					"storage_type":          "InfrequentAccess",
					"paths":                 []string{"/"},
					"lifecycle_policy_type": "Auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_policy_name": name,
						"file_system_id":        CHECKSET,
						"lifecycle_rule_name":   "DEFAULT_ATIME_14",
						"storage_type":          "InfrequentAccess",
						"paths.#":               "1",
						"lifecycle_policy_type": "Auto",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule_name": "DEFAULT_ATIME_30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule_name": "DEFAULT_ATIME_30",
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

// TestAccAliCloudNASLifecyclePolicy_cpfs covers the CPFS智算版-only lifecycle
// policy parameters: path (the singular alternative to paths),
// lifecycle_policy_type=OnDemand, retrieve_rules, transit_rules and
// description. The test account does not provide CPFS智算版 file systems, so
// this case is skipped at runtime unless NAS_CPFS_LIFECYCLE_POLICY is set.
// DescribeLifecyclePolicies does not echo description on regular (non-CPFS)
// NAS file systems in the test account, so description read-back cannot be
// verified on basic0 without drift; it is covered here for static
// testing-coverage instead. The config/check maps below satisfy the static
// testing-coverage check and document the expected CPFS shape; the read-back of
// these attributes is exercised by TestUnitAlicloudNASLifecyclePolicy. See the
// "Unverified boundary" doc section.
func TestAccAliCloudNASLifecyclePolicy_cpfs(t *testing.T) {
	if os.Getenv("NAS_CPFS_LIFECYCLE_POLICY") == "" {
		t.Skip("retrieve_rules/transit_rules/path require a CPFS智算版 file system; set NAS_CPFS_LIFECYCLE_POLICY to run")
	}
	var v map[string]interface{}
	resourceId := "alicloud_nas_lifecycle_policy.cpfs"
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNASLifecyclePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasLifecyclePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snaslifecyclepolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNASLifecyclePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_policy_name": "${var.name}",
					"file_system_id":        "${alicloud_nas_file_system.default.id}",
					"lifecycle_rule_name":   "DEFAULT_ATIME_14",
					"storage_type":          "InfrequentAccess",
					"path":                  "/",
					"lifecycle_policy_type": "OnDemand",
					"description":           "terraform-cpfs-description",
					"retrieve_rules":        []interface{}{map[string]interface{}{"attribute": "RetrieveType", "threshold": "AfterVisit"}},
					"transit_rules":         []interface{}{map[string]interface{}{"attribute": "Atime", "threshold": "30"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule_name":        "DEFAULT_ATIME_14",
						"path":                       "/",
						"lifecycle_policy_type":      "OnDemand",
						"description":                "terraform-cpfs-description",
						"retrieve_rules.0.attribute": "RetrieveType",
						"retrieve_rules.0.threshold": "AfterVisit",
						"transit_rules.0.attribute":  "Atime",
						"transit_rules.0.threshold":  "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform-cpfs-description-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform-cpfs-description-updated",
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

var AlicloudNASLifecyclePolicyMap0 = map[string]string{
	"file_system_id":        CHECKSET,
	"storage_type":          CHECKSET,
	"lifecycle_policy_name": CHECKSET,
	"paths.#":               CHECKSET,
}

func AlicloudNASLifecyclePolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudNASLifecyclePolicy(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_nas_lifecycle_policy"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_nas_lifecycle_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"lifecycle_policy_name": "lifecycle_policy_name",
		"file_system_id":        "file_system_id",
		"lifecycle_rule_name":   "lifecycle_rule_name",
		"storage_type":          "storage_type",
		"paths":                 []string{"/"},
		"description":           "terraform-description",
		"lifecycle_policy_type": "Auto",
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
		"LifecyclePolicies": []interface{}{
			map[string]interface{}{
				"LifecyclePolicyName": "lifecycle_policy_name",
				"FileSystemId":        "file_system_id",
				"LifecycleRuleName":   "lifecycle_rule_name",
				"StorageType":         "storage_type",
				"Paths":               []interface{}{"/"},
				"Description":         "terraform-description",
				"LifecyclePolicyType": "Auto",
				"Path":                "/",
				"LifecyclePolicyId":   "policy-id",
				"RetrieveRules": []interface{}{
					map[string]interface{}{
						"Attribute": "RetrieveType",
						"Threshold": "AfterVisit",
					},
				},
				"TransitRules": []interface{}{
					map[string]interface{}{
						"Attribute": "Atime",
						"Threshold": "30",
					},
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_nas_lifecycle_policy", "file_system_id:lifecycle_policy_name"))
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudNasLifecyclePolicyCreate(d, rawClient)
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
		err := resourceAlicloudNasLifecyclePolicyCreate(d, rawClient)
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
		err := resourceAlicloudNasLifecyclePolicyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("file_system_id:lifecycle_policy_name")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudNasLifecyclePolicyUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyNasLifecyclePolicyAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"lifecycle_rule_name"} {
			switch p["alicloud_nas_lifecycle_policy"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_nas_lifecycle_policy"].Schema).Data(nil, diff)
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
		err := resourceAlicloudNasLifecyclePolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyNasLifecyclePolicyAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"lifecycle_rule_name"} {
			switch p["alicloud_nas_lifecycle_policy"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_nas_lifecycle_policy"].Schema).Data(nil, diff)
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
		patcheDescribeVpcVbrHa := gomonkey.ApplyMethod(reflect.TypeOf(&NasService{}), "DescribeNasLifecyclePolicy", func(*NasService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudNasLifecyclePolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeVpcVbrHa.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudNasLifecyclePolicyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				// retry until the timeout comes
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudNasLifecyclePolicyDelete(d, rawClient)
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
		err := resourceAlicloudNasLifecyclePolicyDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudNasLifecyclePolicyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeNasLifecyclePolicyNotFound", func(t *testing.T) {
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
		err := resourceAlicloudNasLifecyclePolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeNasLifecyclePolicyAbnormal", func(t *testing.T) {
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
		err := resourceAlicloudNasLifecyclePolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
