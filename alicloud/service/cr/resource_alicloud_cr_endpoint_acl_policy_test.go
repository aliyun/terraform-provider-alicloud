package cr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudCREndpointAclPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_endpoint_acl_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCREndpointAclPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrEndpointAclPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	// CR EE instance names are validated server-side (BssOpenApi CreateInstance
	// for the ACR product) and reject uppercase letters with INSTANCE_NAME_INVALID;
	// keep the name lowercase and hyphen-separated to match the convention used
	// by the other alicloud_cr_ee_instance acc tests. The same value is reused
	// as the acl policy description, which has no such constraint.
	name := fmt.Sprintf("tf-testacc-cr-aclep-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCREndpointAclPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"instance_id":   acctest.CHECKSET,
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

// TestAccAliCloudCREndpointAclPolicy_basic1 covers the API-documented
// endpoint_type spelling "Internet": the schema accepts it and normalizes it
// to the lowercase form used in state and resource IDs, and an import ID
// written with "Internet" (as the CR API documentation shows) imports cleanly
// without a perpetual diff.
func TestAccAliCloudCREndpointAclPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_endpoint_acl_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCREndpointAclPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrEndpointAclPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cr-aclep-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCREndpointAclPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   name,
					"entry":         "192.168.2.0/24",
					"instance_id":   "${local.cr_endpoint_instance_id}",
					"module_name":   "Registry",
					"endpoint_type": "Internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"entry":         "192.168.2.0/24",
						"instance_id":   acctest.CHECKSET,
						"module_name":   "Registry",
						"endpoint_type": "internet",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Rewrite the import ID to the API-documented endpoint_type
				// casing; Read must normalize it back to lowercase.
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceId)
					}
					return strings.Replace(rs.Primary.ID, ":internet:", ":Internet:", 1), nil
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
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-loopback-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckCrLoopbackAclPolicyDestroy,
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
// rac.CheckResourceDestroy(), it tolerates INSTANCE_NOT_EXIST from
// GetInstanceEndpoint (via DescribeCrEndpointAclPolicy): this test destroys
// the parent CR EE instance within the same stack, so once the instance is
// deleted the endpoint — and every ACL entry on it — cannot exist anymore.
// Any other error, or a living entry on a living instance, still fails the
// check so real dangling resources are never masked.
func testAccCheckCrLoopbackAclPolicyDestroy(s *terraform.State) error {
	service := crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cr_endpoint_acl_policy" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}
		_, err := service.DescribeCrEndpointAclPolicy(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) || strings.Contains(err.Error(), "INSTANCE_NOT_EXIST") {
				continue
			}
			return errs.WrapError(err)
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
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
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
func TestUnitAlicloudCrEndpointTypeConvert(t *testing.T) {
	for input, want := range map[string]string{
		"internet": "Internet",
		"Internet": "Internet",
		"INTERNET": "Internet",
		"vpc":      "vpc",
		"":         "",
	} {
		assert.Equal(t, want, crsvc.ConvertCrEndpointType(input), "crsvc.ConvertCrEndpointType(%q)", input)
	}
	assert.True(t, crsvc.IsCrEndpointType("internet"))
	assert.True(t, crsvc.IsCrEndpointType("Internet"))
	assert.False(t, crsvc.IsCrEndpointType("10.0.0.0/8"))
	assert.False(t, crsvc.IsCrEndpointType("vpc"))
	assert.False(t, crsvc.IsCrEndpointType(""))
}

// TestUnitAlicloudCREndpointAclPolicyReadInvalidId reproduces the misordered
// import ID from the field report (<instance_id>:<entry>:<endpoint_type>
// instead of the documented <instance_id>:<endpoint_type>:<entry>): Read must
// fail fast with an actionable message naming the ID format instead of
// forwarding the segment to the API and surfacing a raw InvalidEndpointType.
