package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudEbsEnterpriseSnapshotPolicy_basic2269(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap2269)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEbsEnterpriseSnapshotPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence2269)
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
					"status": "ENABLED",
					"desc":   "${var.name}",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 */12 * * *",
						},
					},
					"storage_rule": []map[string]interface{}{
						{
							"enable_immediate_access": "false",
						},
					},
					"target_type": "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "1",
							"time_unit":     "DAYS",
						},
					},
					"cross_region_copy_info": []map[string]interface{}{
						{
							"enabled": "true",
							"regions": []map[string]interface{}{
								{
									"region_id":   defaultRegionToTest,
									"retain_days": "7",
								},
							},
						},
					},
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "Test1",
					},
					"enterprise_snapshot_policy_name": "${var.name}",
					"resource_group_id":               "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "ENABLED",
						"desc":                            name,
						"target_type":                     "DISK",
						"enterprise_snapshot_policy_name": name,
						"resource_group_id":               CHECKSET,
						"schedule.#":                      "1",
						"schedule.0.cron_expression":      "0 0 */12 * * *",
						"tags.%":                          "2",
						"tags.Created":                    "TF1",
						"tags.For":                        "Test1",
						"retain_rule.#":                   "1",
						"retain_rule.0.time_interval":     "1",
						"retain_rule.0.time_unit":         "DAYS",
						"storage_rule.#":                  "1",
						"cross_region_copy_info.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLED",
					"desc":   "${var.name}_update",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 */6 * * *",
						},
					},
					"retain_rule": []map[string]interface{}{
						{
							"number":        "1",
							"time_interval": "1",
							"time_unit":     "DAYS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                      "DISABLED",
						"desc":                        name + "_update",
						"schedule.#":                  "1",
						"schedule.0.cron_expression":  "0 0 */6 * * *",
						"retain_rule.#":               "1",
						"retain_rule.0.number":        "1",
						"retain_rule.0.time_interval": "1",
						"retain_rule.0.time_unit":     "DAYS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "Test2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF2",
						"tags.For":     "Test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enterprise_snapshot_policy_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_region_copy_info": []map[string]interface{}{
						{
							"enabled": "true",
							"regions": []map[string]interface{}{
								{
									"region_id":   defaultRegionToTest,
									"retain_days": "30",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_region_copy_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_rule": []map[string]interface{}{
						{
							"enable_immediate_access": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_rule.#": "1",
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

var AlicloudEbsEnterpriseSnapshotPolicyMap2269 = map[string]string{}

func AlicloudEbsEnterpriseSnapshotPolicyBasicDependence2269(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
