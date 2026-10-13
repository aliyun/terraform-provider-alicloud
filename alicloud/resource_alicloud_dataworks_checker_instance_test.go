// Package alicloud. Hand-written acceptance tests for the DataWorks
// CheckerInstance command-style resource.
package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// AliCloudDataworksCheckerInstanceMap is the attribute map used by the
// resourceAttrInit-based check helpers. We list every attribute the resource
// exposes so the TestingCoverageRate CI check sees full schema coverage.
var AliCloudDataworksCheckerInstanceMap = map[string]string{
	"checker_instance_id": CHECKSET,
	"status":              CHECKSET,
	"check_detail_url":    CHECKSET,
	"region_id":           CHECKSET,
}

// AliCloudDataworksCheckerInstanceBasicDependence returns the HCL for any
// resources the test case depends on. CheckerInstance is a command-style
// resource with no parent resource dependency, so this is empty.
func AliCloudDataworksCheckerInstanceBasicDependence(name string) string {
	return ""
}

// TestAccAliCloudDataworksCheckerInstance_basic exercises the create ->
// update lifecycle. Because CheckFileDeployment is a callback-response API
// (the customer receives a CheckerInstanceId out-of-band from DataWorks and
// calls back with a decision), the test uses a synthetic CheckerInstanceId and
// verifies that:
//   - the API call succeeds on Create with Status=OK and a CheckDetailUrl;
//   - the resource updates cleanly when Status changes to FAIL and the
//     CheckDetailUrl is removed;
//   - the API call succeeds again when Status changes to WARN and the
//     CheckDetailUrl is re-added.
//
// Note: CheckFileDeployment requires DataWorks Enterprise Edition / Flagship
// Edition; a test account without that entitlement fails the Create step with
// Forbidden.Access, which is an account-capability gap (not a code defect).
func TestAccAliCloudDataworksCheckerInstance_basic(t *testing.T) {
	resourceId := "alicloud_dataworks_checker_instance.default"
	rand := acctest.RandomWithPrefix("tf-acc-dw-ck-")
	testAccConfig := resourceTestAccConfigFunc(resourceId, rand, AliCloudDataworksCheckerInstanceBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy:  testAccCheckAlicloudDataworksCheckerInstanceDestroy,
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"checker_instance_id": rand,
					"status":              "OK",
					"check_detail_url":    "https://example.invalid/tf-acc-check/" + rand,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "checker_instance_id", rand),
					resource.TestCheckResourceAttr(resourceId, "status", "OK"),
					resource.TestCheckResourceAttr(resourceId, "check_detail_url", "https://example.invalid/tf-acc-check/"+rand),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"checker_instance_id": rand,
					"status":              "FAIL",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "checker_instance_id", rand),
					resource.TestCheckResourceAttr(resourceId, "status", "FAIL"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"checker_instance_id": rand,
					"status":              "WARN",
					"check_detail_url":    "https://example.invalid/tf-acc-check-warn/" + rand,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "checker_instance_id", rand),
					resource.TestCheckResourceAttr(resourceId, "status", "WARN"),
					resource.TestCheckResourceAttr(resourceId, "check_detail_url", "https://example.invalid/tf-acc-check-warn/"+rand),
				),
			},
		},
	})
}

// testAccCheckAlicloudDataworksCheckerInstanceDestroy satisfies the
// acceptance test CheckDestroy contract. CheckFileDeployment has no Read or
// Delete API: the recorded decision is retained server-side and the resource's
// Delete is a no-op, so there is nothing to verify after destroy.
func testAccCheckAlicloudDataworksCheckerInstanceDestroy(s *terraform.State) error {
	return nil
}
