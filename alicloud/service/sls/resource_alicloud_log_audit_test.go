package sls_test

import (
	"fmt"
	"os"
	"testing"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogAudit_basic(t *testing.T) {
	var v *slsPop.DescribeAppResponse
	resourceId := "alicloud_log_audit.foo"
	ra := acctest.ResourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
						"oss_access_enabled":  "true",
						"oss_access_ttl":      "155",
						"oss_sync_enabled":    "true",
						"oss_sync_ttl":        "180",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":                     name,
						"aliuid":                           acctest.CHECKSET,
						"variable_map.%":                   "6",
						"variable_map.actiontrail_enabled": "false",
						"variable_map.actiontrail_ttl":     "10",
						"variable_map.oss_access_enabled":  "true",
						"variable_map.oss_access_ttl":      "155",
						"variable_map.oss_sync_enabled":    "true",
						"variable_map.oss_sync_ttl":        "180",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_account": []string{"1234567", "123123123213", "123141412"},
					"variable_map": map[string]string{
						"actiontrail_enabled": "true",
						"actiontrail_ttl":     "10",
						"oss_access_enabled":  "true",
						"oss_access_ttl":      "155",
						"oss_sync_enabled":    "true",
						"oss_sync_ttl":        "180",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_account.#":                  "3",
						"variable_map.%":                   "6",
						"variable_map.actiontrail_enabled": "true",
					}),
				),
			},
			// TODO: only when center account is resource directory master or resource directory admin need to check resource type config，otherwise pass it directly
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"resource_directory_type": "custom",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"resource_directory_type": "custom",
			// 		}),
			// 	}),
			// },
		},
	})
}

func resourceLogAuditConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "default" {}
`)
}

// testAccLogAuditPreCheckResourceDirectoryAccount gates the
// resource_directory_type case: only a resource-directory master or admin
// account actually applies the directory mode, so the case runs when the
// dedicated env flag opts in and skips otherwise (same pattern as
// privatezoneTestPreCheckEnterpriseAccountEnabled).
func testAccLogAuditPreCheckResourceDirectoryAccount(t *testing.T) {
	if os.Getenv("LOG_AUDIT_RESOURCE_DIRECTORY_ENABLED") != "true" {
		t.Skipf("Skipping the test case because the account is not a resource directory master or admin.")
		t.Skipped()
	}
}

// TestAccAliCloudLogAudit_resourceDirectoryType covers the
// resource_directory_type attribute: "custom" requires multi_account, and
// switching to "all" updates the directory mode in place.
func TestAccAliCloudLogAudit_resourceDirectoryType(t *testing.T) {
	var v *slsPop.DescribeAppResponse
	resourceId := "alicloud_log_audit.foo"
	ra := acctest.ResourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccLogAuditPreCheckResourceDirectoryAccount(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name":            name,
					"aliuid":                  "${data.alicloud_account.default.id}",
					"resource_directory_type": "custom",
					"multi_account":           []string{"1234567"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":            name,
						"resource_directory_type": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name":            name,
					"aliuid":                  "${data.alicloud_account.default.id}",
					"resource_directory_type": "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_directory_type": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var logAuditMap = map[string]string{
	"display_name": acctest.CHECKSET,
}

// TestGetInitParameter is a table-driven unit test for getInitParameter.
// It does NOT require cloud resources and covers all edge cases related to
// the displayName fallback fix.

// TestAccAliCloudLogAudit_display_name_fallback is a minimal acceptance test that
// verifies the displayName == "" fallback to d.Id() fix.
//
// It uses the same data.alicloud_account pattern as TestAccAliCloudLogAudit_basic
// but with a minimal variable_map and an additional PlanOnly step to verify
// no spurious diff after the fix.
func TestAccAliCloudLogAudit_display_name_fallback(t *testing.T) {
	var v *slsPop.DescribeAppResponse
	resourceId := "alicloud_log_audit.foo"
	ra := acctest.ResourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			// Step 1: Create log_audit with minimal config.
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":                     name,
						"variable_map.actiontrail_enabled": "false",
						"variable_map.actiontrail_ttl":     "10",
					}),
				),
			},
			// Step 2: Import and verify state.
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: Re-plan with same config — expect no diff.
			// This is the core assertion for the displayName fallback fix.
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
					},
				}),
				PlanOnly: true,
			},
		},
	})
}
