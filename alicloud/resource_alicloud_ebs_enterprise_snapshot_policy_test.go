package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ebs EnterpriseSnapshotPolicy. >>> Resource test cases, automatically generated.
// Case 5473
func TestAccAliCloudEbsEnterpriseSnapshotPolicy_basic5473(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap5473)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5473)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name,
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "1",
							"time_unit":     "WEEKS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name,
						"target_type":                     "DISK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ENABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desc": "ESP 资源测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desc": "ESP 资源测试",
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
									"region_id":   "cn-hangzhou-test-4",
									"retain_days": "14",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enterprise_snapshot_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "2",
							"time_unit":     "WEEKS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desc": "NewDesc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desc": "NewDesc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 2 2 * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enterprise_snapshot_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_rule": []map[string]interface{}{
						{
							"time_unit": "MONTHS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_region_copy_info": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 2 3 * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 */8 * * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
					"desc":   "ESP 资源测试",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name + "_update",
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "1",
							"time_unit":     "WEEKS",
						},
					},
					"cross_region_copy_info": []map[string]interface{}{
						{
							"enabled": "true",
							"regions": []map[string]interface{}{
								{
									"region_id":   "cn-hangzhou-test-4",
									"retain_days": "14",
								},
							},
						},
					},
					"special_retain_rules": []map[string]interface{}{
						{
							"enabled": "true",
							"rules": []map[string]interface{}{
								{
									"special_period_unit": "WEEKS",
									"time_interval":       "16",
									"time_unit":           "WEEKS",
								},
								{
									"special_period_unit": "MONTHS",
									"time_interval":       "1",
									"time_unit":           "WEEKS",
								},
								{
									"special_period_unit": "WEEKS",
									"time_interval":       "12",
									"time_unit":           "WEEKS",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "ENABLED",
						"desc":                            "ESP 资源测试",
						"enterprise_snapshot_policy_name": name + "_update",
						"target_type":                     "DISK",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEbsEnterpriseSnapshotPolicyMap5473 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5473(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5494
func TestAccAliCloudEbsEnterpriseSnapshotPolicy_basic5494(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap5494)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5494)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name,
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "120",
							"time_unit":     "DAYS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name,
						"target_type":                     "DISK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desc": "DESC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desc": "DESC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enterprise_snapshot_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "120",
							"time_unit":     "DAYS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "DISABLED",
					"desc":              "DESC",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name + "_update",
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "120",
							"time_unit":     "DAYS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "DISABLED",
						"desc":                            "DESC",
						"resource_group_id":               CHECKSET,
						"enterprise_snapshot_policy_name": name + "_update",
						"target_type":                     "DISK",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEbsEnterpriseSnapshotPolicyMap5494 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5494(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)
}

// Case 5484
func TestAccAliCloudEbsEnterpriseSnapshotPolicy_basic5484(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap5484)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5484)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name,
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"number": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name,
						"target_type":                     "DISK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ENABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desc": "DESC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desc": "DESC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enterprise_snapshot_policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enterprise_snapshot_policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_rule": []map[string]interface{}{
						{
							"number": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 */24 * * ?",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retain_rule": []map[string]interface{}{
						{
							"number": "12",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
					"desc":   "DESC",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name + "_update",
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"number": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "ENABLED",
						"desc":                            "DESC",
						"enterprise_snapshot_policy_name": name + "_update",
						"target_type":                     "DISK",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEbsEnterpriseSnapshotPolicyMap5484 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5484(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5473  twin
func TestAccAliCloudEbsEnterpriseSnapshotPolicy_basic5473_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap5473)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5473)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLED",
					"desc":   "NewDesc",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 */8 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name,
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "1",
							"time_unit":     "MONTHS",
						},
					},
					"storage_rule": []map[string]interface{}{
						{
							"enable_immediate_access": "true",
						},
					},
					"cross_region_copy_info": []map[string]interface{}{
						{
							"enabled": "true",
							"regions": []map[string]interface{}{
								{
									"region_id":   "cn-hangzhou-test6",
									"retain_days": "129",
								},
							},
						},
					},
					"special_retain_rules": []map[string]interface{}{
						{
							"enabled": "true",
							"rules": []map[string]interface{}{
								{
									"special_period_unit": "MONTHS",
									"time_interval":       "4",
									"time_unit":           "YEARS",
								},
								{
									"special_period_unit": "MONTHS",
									"time_interval":       "1",
									"time_unit":           "WEEKS",
								},
								{
									"special_period_unit": "WEEKS",
									"time_interval":       "12",
									"time_unit":           "WEEKS",
								},
							},
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "DISABLED",
						"desc":                            "NewDesc",
						"enterprise_snapshot_policy_name": name,
						"target_type":                     "DISK",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
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

// Case 5494  twin
func TestAccAliCloudEbsEnterpriseSnapshotPolicy_basic5494_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap5494)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5494)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "DISABLED",
					"desc":              "DESC",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 1 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name,
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"time_interval": "120",
							"time_unit":     "DAYS",
							"number":        "12",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "DISABLED",
						"desc":                            "DESC",
						"resource_group_id":               CHECKSET,
						"enterprise_snapshot_policy_name": name,
						"target_type":                     "DISK",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
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

// Case 5484  twin
func TestAccAliCloudEbsEnterpriseSnapshotPolicy_basic5484_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_enterprise_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsEnterpriseSnapshotPolicyMap5484)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsEnterpriseSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebsenterprisesnapshotpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsEnterpriseSnapshotPolicyBasicDependence5484)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EbsEnterpriseSnapshotPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
					"desc":   "DESC",
					"schedule": []map[string]interface{}{
						{
							"cron_expression": "0 0 */24 * * ?",
						},
					},
					"enterprise_snapshot_policy_name": name,
					"target_type":                     "DISK",
					"retain_rule": []map[string]interface{}{
						{
							"number": "1",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                          "ENABLED",
						"desc":                            "DESC",
						"enterprise_snapshot_policy_name": name,
						"target_type":                     "DISK",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
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

// Test Ebs EnterpriseSnapshotPolicy. <<< Resource test cases, automatically generated.
