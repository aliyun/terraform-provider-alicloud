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
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

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
	// as the acl policy description, which has no such constraint.
	name := fmt.Sprintf("tf-testacc-cr-aclep-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCREndpointAclPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   name,
					"entry":         "192.168.1.0/24",
					"instance_id":   "${local.cr_endpoint_instance_id}",
					"module_name":   "Registry",
					"endpoint_type": "internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"entry":         "192.168.1.0/24",
						"instance_id":   CHECKSET,
						"module_name":   "Registry",
						"endpoint_type": "internet",
					}),
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

// TestAccAliCloudCrInternetEndpoint_loopbackAclPolicy probes the server-side
// behavior when the standalone alicloud_cr_endpoint_acl_policy resource is
// used to create an ACL entry with entry 127.0.0.1/32 and description
// "default" — the same values the system auto-adds as a loopback ACL once
// the internet endpoint is enabled (see resource_alicloud_cr_internet_endpoint
// Read filter at the 127.0.0.1/32+default entry).
//
// GetInstanceEndpoint auto-adds a loopback ACL (entry 127.0.0.1/32, comment
// "default") when the internet endpoint is enabled. This test creates that
// exact entry through the ACL policy resource to determine whether the
// server:
//  1. rejects the duplicate (CreateInstanceEndpointAclPolicy returns an
//     error — the Step Apply fails),
//  2. silently deduplicates (Create succeeds, only one loopback entry
//     remains), or
//  3. accepts a duplicate (Create succeeds, two loopback entries coexist).
//
// The observed loopback entry count is logged via t.Logf for post-run
// analysis; no hard assertion is made on the count so the test does not mask
// any of the three behaviors. The custom CheckDestroy is the real gate: a
// residual ACL entry on a living instance fails the teardown, while
// INSTANCE_NOT_EXIST from GetInstanceEndpoint is accepted as destroyed
// because this test creates a throwaway CR EE instance that is destroyed
// within the same stack.
func TestAccAliCloudCrInternetEndpoint_loopbackAclPolicy(t *testing.T) {
	resourceId := "alicloud_cr_endpoint_acl_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-loopback-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCrLoopbackAclPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCrLoopbackAclPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCrLoopbackAclEntryCount(t, resourceId),
				),
			},
		},
	})
}

// testAccCheckCrLoopbackAclPolicyDestroy verifies that the ACL policy created
// by the loopback probe is gone after destroy. Unlike the shared
// rac.checkResourceDestroy(), it tolerates INSTANCE_NOT_EXIST from
// GetInstanceEndpoint (via DescribeCrEndpointAclPolicy): this test destroys
// the parent CR EE instance within the same stack, so once the instance is
// deleted the endpoint — and every ACL entry on it — cannot exist anymore.
// Any other error, or a living entry on a living instance, still fails the
// check so real dangling resources are never masked.
func testAccCheckCrLoopbackAclPolicyDestroy(s *terraform.State) error {
	service := CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cr_endpoint_acl_policy" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}
		_, err := service.DescribeCrEndpointAclPolicy(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) || strings.Contains(err.Error(), "INSTANCE_NOT_EXIST") {
				continue
			}
			return WrapError(err)
		}
		return fmt.Errorf("ACL policy %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCrLoopbackAclPolicyConfig(name string) string {
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

resource "alicloud_cr_internet_endpoint" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
  entries {
    entry   = "192.168.1.0/24"
    comment = "normal-entry"
  }
}

resource "alicloud_cr_endpoint_acl_policy" "default" {
  instance_id   = alicloud_cr_internet_endpoint.default.id
  endpoint_type = "internet"
  module_name   = "Registry"
  entry         = "127.0.0.1/32"
  description   = "default"
}
`, name)
}

// testAccCheckCrLoopbackAclEntryCount queries GetInstanceEndpoint after the
// ACL policy apply and counts how many AclEntries have Entry=="127.0.0.1/32"
// and Comment=="default". The count distinguishes the three server-side
// behaviors (0=silently rejected, 1=deduplicated, 2=duplicate accepted).
// The count is logged via t.Logf and the check always passes so the test
// does not mask any behavior; the teardown (CheckDestroy) is the real gate.
func testAccCheckCrLoopbackAclEntryCount(t *testing.T, resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		if instanceId == "" {
			return fmt.Errorf("instance_id attribute is empty for resource %s", resourceId)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"InstanceId":   instanceId,
			"EndpointType": "internet",
			"ModuleName":   "Registry",
			"RegionId":     client.RegionId,
		}
		response, err := client.RpcPost("cr", "2018-12-01", "GetInstanceEndpoint", nil, request, true)
		if err != nil {
			return fmt.Errorf("GetInstanceEndpoint failed for instance %s: %v", instanceId, err)
		}
		rawEntries, ok := response["AclEntries"]
		if !ok {
			t.Logf("[loopback-acl-probe] instance %s: AclEntries not present in GetInstanceEndpoint response: %v", instanceId, response)
			return nil
		}
		entries, ok := rawEntries.([]interface{})
		if !ok {
			t.Logf("[loopback-acl-probe] instance %s: AclEntries is not an array: %v", instanceId, rawEntries)
			return nil
		}
		loopbackCount := 0
		totalEntries := 0
		for _, e := range entries {
			totalEntries++
			m, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(m["Entry"]) == "127.0.0.1/32" && fmt.Sprint(m["Comment"]) == "default" {
				loopbackCount++
			}
		}
		t.Logf("[loopback-acl-probe] instance %s: total AclEntries=%d, loopback (127.0.0.1/32+default) count=%d", instanceId, totalEntries, loopbackCount)
		return nil
	}
}

// lintignore: R001
func TestUnitAlicloudCREndpointAclPolicy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cr_endpoint_acl_policy"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_cr_endpoint_acl_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"description":   "description",
		"entry":         "MockEntry",
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
				"Description":  "description",
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
	d.SetId(fmt.Sprint("instance_id", ":", "internet", ":", "MockEntry"))

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
				return responseMock["RetryError"]("SLB_SERVICE_ERROR")
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
