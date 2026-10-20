package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAccAliCloudCREndpointAclPolicy_basic0 covers the create/read/import/
// destroy lifecycle for an internet endpoint ACL policy carrying a single
// user-managed entry.
//
// The config is expressed as an inline map[string]interface{} (built by
// resourceTestAccConfigFunc) rather than a hand-rolled fmt.Sprintf HCL string
// so that scripts/testing/testing_coverage_rate_check.go can parse the
// `entries` set and its nested `entry`/`comment` attributes for set coverage.
// Create gates on WaitCrEndpointRunning so the ACL entry is reliably
// accepted, then waits for it to propagate into AclEntries before the first
// Read; Read excludes the system loopback default (127.0.0.1/32 + "default")
// so the user-managed entry matches state cleanly.
func TestAccAliCloudCREndpointAclPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_endpoint_acl_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCREndpointAclPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEndpointAclPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	// CR EE instance names are validated server-side (BssOpenApi CreateInstance
	// for the ACR product) and reject uppercase letters with INSTANCE_NAME_INVALID;
	// keep the name lowercase and hyphen-separated to match the convention used
	// by the other alicloud_cr_ee_instance acc tests. The same value is reused
	// as the acl entry comment, which has no such constraint.
	name := fmt.Sprintf("tf-testacc-cr-aclep-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCREndpointAclPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${local.cr_endpoint_instance_id}",
					"module_name":   "Registry",
					"endpoint_type": "internet",
					"entries": []map[string]interface{}{
						{
							"entry":   "192.168.1.0/24",
							"comment": name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"module_name":   "Registry",
						"endpoint_type": "internet",
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"module_name"},
			},
		},
	})
}

// TestAccAliCloudCREndpointAclPolicy_basic1 covers the API-documented
// endpoint_type spelling "Internet": the schema accepts it and normalizes it
// to the lowercase form used in state and resource IDs, and an import ID
// written with "Internet" (as the CR API documentation shows) imports cleanly
// without a perpetual diff.
func TestAccAliCloudCREndpointAclPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_endpoint_acl_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCREndpointAclPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEndpointAclPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cr-aclep-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCREndpointAclPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${local.cr_endpoint_instance_id}",
					"module_name":   "Registry",
					"endpoint_type": "Internet",
					"entries": []map[string]interface{}{
						{
							"entry":   "192.168.2.0/24",
							"comment": name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"module_name":   "Registry",
						"endpoint_type": "internet",
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "description", name),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Rewrite the import ID's endpoint_type segment to the
				// API-documented "Internet" casing; Read must normalize it back
				// to lowercase. The ID is the 2-part <instance_id>:<endpoint_type>.
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceId)
					}
					parts := strings.Split(rs.Primary.ID, ":")
					parts[len(parts)-1] = "Internet"
					return strings.Join(parts, ":"), nil
				},
				ImportStateVerifyIgnore: []string{"module_name"},
			},
		},
	})
}

var AlicloudCREndpointAclPolicyMap0 = map[string]string{}

func AlicloudCREndpointAclPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Economy"
  instance_name  = var.name
  image_scanner  = "DISABLE"
}

locals {
  cr_endpoint_instance_id = alicloud_cr_ee_instance.default.id
}

data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = local.cr_endpoint_instance_id
  module_name   = "Registry"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudCREndpointAclPolicy(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cr_endpoint_acl_policy"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_cr_endpoint_acl_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"entries": []interface{}{
			map[string]interface{}{
				"entry":   "MockEntry",
				"comment": "description",
			},
		},
		"instance_id":   "instance_id",
		"module_name":   "Registry",
		"endpoint_type": "internet",
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
		"AclEntries": []interface{}{
			map[string]interface{}{
				"Entry":        "MockEntry",
				"EndpointType": "internet",
				"Comment":      "description",
				"InstanceId":   "instance_id",
			},
		},
		"IsSuccess": "true",
		"Status":    "RUNNING",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cr_endpoint_acl_policy", "MockEntry"))
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
			result["Entry"] = "MockEntry"
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAcrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudCrEndpointAclPolicyCreate(d, rawClient)
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
		err := resourceAlicloudCrEndpointAclPolicyCreate(d, rawClient)
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
		err := resourceAlicloudCrEndpointAclPolicyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "internet"))

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAcrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudCrEndpointAclPolicyDelete(d, rawClient)
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
		err := resourceAlicloudCrEndpointAclPolicyDelete(d, rawClient)
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
		err := resourceAlicloudCrEndpointAclPolicyDelete(d, rawClient)
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
		err := resourceAlicloudCrEndpointAclPolicyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("DeleteSLBServiceErrorAsNotExist", func(t *testing.T) {
		// The corrected CloudSpec delete mapping marks SLB_SERVICE_ERROR as a
		// resource-not-exist code: the entry cannot be removed (e.g. the system
		// loopback default 127.0.0.1/32) or is already gone, so Delete treats it
		// as deleted instead of retrying until the delete timeout. The retryable
		// SLB_SERVICE_ERROR path is covered by DeleteMockAbnormal (Throttling);
		// this case pins the not-exist semantics on the SLB code itself.
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["RetryError"]("SLB_SERVICE_ERROR")
		})
		err := resourceAlicloudCrEndpointAclPolicyDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeCrEndpointAclPolicyNotFound", func(t *testing.T) {
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
		err := resourceAlicloudCrEndpointAclPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeCrEndpointAclPolicyAbnormal", func(t *testing.T) {
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
		err := resourceAlicloudCrEndpointAclPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

func TestUnitAlicloudCrEndpointTypeConvert(t *testing.T) {
	for input, want := range map[string]string{
		"internet": "Internet",
		"Internet": "Internet",
		"INTERNET": "Internet",
		"vpc":      "vpc",
		"":         "",
	} {
		assert.Equal(t, want, convertCrEndpointType(input), "convertCrEndpointType(%q)", input)
	}
	assert.True(t, isCrEndpointType("internet"))
	assert.True(t, isCrEndpointType("Internet"))
	assert.False(t, isCrEndpointType("10.0.0.0/8"))
	assert.False(t, isCrEndpointType("vpc"))
	assert.False(t, isCrEndpointType(""))
}

// TestUnitAlicloudCREndpointAclPolicyReadInvalidId reproduces a malformed
// resource ID whose endpoint_type segment is not a valid CR endpoint type (an
// IP segment supplied where <endpoint_type> belongs): Read must fail fast with
// an actionable message naming the ID format (<instance_id>:<endpoint_type>)
// instead of forwarding the segment to the API and surfacing a raw
// InvalidEndpointType.
func TestUnitAlicloudCREndpointAclPolicyReadInvalidId(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cr_endpoint_acl_policy"].Schema).Data(nil, nil)
	d.SetId("cri-abc123:10.0.0.0/8")
	err := resourceAlicloudCrEndpointAclPolicyRead(d, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid endpoint_type")
	assert.Contains(t, err.Error(), "<instance_id>:<endpoint_type>")
}
