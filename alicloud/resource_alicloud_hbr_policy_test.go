package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Hbr Policy. >>> Resource test cases, automatically generated.
// Case Policy 6287
func TestAccAliCloudHbrPolicy_basic6287(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyMap6287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBasicDependence6287)
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
					"policy_name": name,
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"archive_days": "0",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
						},
						{
							"rule_type":    "TRANSITION",
							"backup_type":  "COMPLETE",
							"retention":    "120",
							"archive_days": "0",
						},
						{
							"rule_type":             "REPLICATION",
							"backup_type":           "COMPLETE",
							"retention":             "135",
							"replication_region_id": "cn-chengdu",
							"archive_days":          "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name,
						"rules.#":     "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "policy creation",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "120",
							"archive_days": "30",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "240",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "0",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "175",
							"replication_region_id": "cn-qingdao",
							"archive_days":          "0",
							"backup_type":           "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "policy update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "240",
							"archive_days": "85",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P2D",
							"retention":             "7",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "120",
							"replication_region_id": "cn-zhangjiakou",
							"archive_days":          "50",
							"backup_type":           "COMPLETE",
						},
						{
							"rule_type":    "BACKUP",
							"backup_type":  "INCREMENTAL",
							"schedule":     "I|1631685600|PT1H",
							"retention":    "7",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"archive_days": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"archive_days": "0",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
						},
						{
							"rule_type":    "TRANSITION",
							"retention":    "120",
							"archive_days": "0",
							"backup_type":  "COMPLETE",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "135",
							"replication_region_id": "cn-chengdu",
							"archive_days":          "0",
							"backup_type":           "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"archive_days": "0",
						},
						{
							"rule_type":    "TRANSITION",
							"retention":    "145",
							"archive_days": "0",
							"backup_type":  "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "policy creation",
					"policy_name":        name + "_update",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "120",
							"archive_days": "30",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "240",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "0",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "175",
							"replication_region_id": "cn-qingdao",
							"archive_days":          "0",
							"backup_type":           "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name + "_update",
						"rules.#":            "3",
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

var AliCloudHbrPolicyMap6287 = map[string]string{
	"create_time": CHECKSET,
	"policy_type": CHECKSET,
}

func AliCloudHbrPolicyBasicDependence6287(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_hbr_vault" "defaulth4dKAG" {
  		vault_type          = "STANDARD"
  		encrypt_type        = "HBR_PRIVATE"
  		vault_name          = var.name
  		vault_storage_class = "STANDARD"
	}

	resource "alicloud_hbr_vault" "defaultL7kwwD" {
  		vault_type          = "STANDARD"
  		encrypt_type        = "HBR_PRIVATE"
  		vault_name          = join("-", [var.name, 1])
  		vault_storage_class = "STANDARD"
	}
`, name)
}

// Case Policy测试用例 5320
func TestAccAliCloudHbrPolicy_basic5320(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyMap5320)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBasicDependence5320)
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
					"policy_name": name,
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"backup_type":  "COMPLETE",
							"retention":    "120",
							"archive_days": "30",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "240",
								},
							},
							"replication_region_id": "cn-shanghai",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "0",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"replication_region_id": "cn-beijing",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name,
						"rules.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "镇元Policy-创建",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元Policy-创建",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "镇元-修改",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元-修改",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "240",
							"archive_days": "85",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"backup_type":           "COMPLETE",
							"replication_region_id": "cn-beijing",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P2D",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
							"replication_region_id": "cn-shanghai",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "8",
							"replication_region_id": "cn-beijing",
							"archive_days":          "50",
							"backup_type":           "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "镇元Policy-创建",
					"policy_name":        name + "_update",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"backup_type":  "COMPLETE",
							"retention":    "120",
							"archive_days": "30",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "240",
								},
							},
							"replication_region_id": "cn-shanghai",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "0",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"replication_region_id": "cn-beijing",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元Policy-创建",
						"policy_name":        name + "_update",
						"rules.#":            "2",
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

var AliCloudHbrPolicyMap5320 = map[string]string{
	"create_time": CHECKSET,
	"policy_type": CHECKSET,
}

func AliCloudHbrPolicyBasicDependence5320(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_hbr_vault" "defaulth4dKAG" {
  		vault_type          = "STANDARD"
  		encrypt_type        = "HBR_PRIVATE"
  		vault_name          = var.name
  		vault_storage_class = "STANDARD"
	}

	resource "alicloud_hbr_vault" "defaultL7kwwD" {
  		vault_type          = "STANDARD"
  		encrypt_type        = "HBR_PRIVATE"
  		vault_name          = join("-", [var.name, 1])
  		vault_storage_class = "STANDARD"
	}
`, name)
}

// Case Policy 6287  twin
func TestAccAliCloudHbrPolicy_basic6287_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyMap6287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBasicDependence6287)
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
					"policy_description": "policy update",
					"policy_name":        name,
					"policy_type":        "STANDARD",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"retention":    "240",
							"archive_days": "0",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"backup_type": "COMPLETE",
							"schedule":    "I|1631685600|P1D",
							"vault_id":    "${alicloud_hbr_vault.defaultL7kwwD.id}",
						},
						{
							"rule_type":             "TRANSITION",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P2D",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"retention":             "145",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "135",
							"replication_region_id": "cn-chengdu",
							"archive_days":          "0",
							"backup_type":           "COMPLETE",
						},
						{
							"rule_type":    "BACKUP",
							"backup_type":  "INCREMENTAL",
							"schedule":     "I|1631685600|PT1H",
							"retention":    "240",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"archive_days": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy update",
						"policy_name":        name,
						"policy_type":        "STANDARD",
						"rules.#":            "4",
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

// Case Policy测试用例 5320  twin
func TestAccAliCloudHbrPolicy_basic5320_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyMap5320)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBasicDependence5320)
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
					"policy_description": "镇元-修改",
					"policy_name":        name,
					"policy_type":        "UDM_ECS_ONLY",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"retention":    "8",
							"archive_days": "0",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"replication_region_id": "cn-chengdu",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "INCREMENTAL",
							"schedule":              "I|1631685600|P2D",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"replication_region_id": "cn-chengdu",
							"retention":             "8",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "8",
							"replication_region_id": "cn-beijing",
							"archive_days":          "50",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"backup_type": "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元-修改",
						"policy_name":        name,
						"policy_type":        "UDM_ECS_ONLY",
						"rules.#":            "3",
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

// The original test case on the api management

// Case Policy 6287   raw
func TestAccAliCloudHbrPolicy_basic6287_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyMap6287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBasicDependence6287)
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
					"policy_description": "policy creation",
					"policy_name":        name,
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "120",
							"archive_days": "30",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "240",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "0",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "175",
							"replication_region_id": "cn-qingdao",
							"archive_days":          "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name,
						"rules.#":            "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "policy update",
					"policy_name":        name + "_update",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "240",
							"archive_days": "85",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"retention":             "7",
							"schedule":              "I|1631685600|P2D",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "120",
							"replication_region_id": "cn-zhangjiakou",
							"archive_days":          "50",
							"backup_type":           "COMPLETE",
						},
						{
							"rule_type":    "BACKUP",
							"backup_type":  "INCREMENTAL",
							"schedule":     "I|1631685600|PT1H",
							"retention":    "7",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"archive_days": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy update",
						"policy_name":        name + "_update",
						"rules.#":            "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "_update",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"archive_days": "0",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
						},
						{
							"rule_type":    "TRANSITION",
							"retention":    "120",
							"archive_days": "0",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "135",
							"replication_region_id": "cn-chengdu",
							"archive_days":          "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
						"rules.#":     "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"archive_days": "0",
						},
						{
							"rule_type":    "TRANSITION",
							"retention":    "145",
							"archive_days": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rules.#": "2",
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

// Case Policy测试用例 5320   raw
func TestAccAliCloudHbrPolicy_basic5320_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyMap5320)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBasicDependence5320)
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
					"policy_description": "镇元Policy-创建",
					"policy_name":        name,
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "120",
							"archive_days": "30",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "240",
								},
							},
							"replication_region_id": "cn-shanghai",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "0",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"replication_region_id": "cn-beijing",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元Policy-创建",
						"policy_name":        name,
						"rules.#":            "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_description": "镇元-修改",
					"policy_name":        name + "_update",
					"rules": []map[string]interface{}{
						{
							"rule_type":    "TRANSITION",
							"retention":    "240",
							"archive_days": "85",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "WEEKLY",
									"retention":               "480",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "960",
								},
								{
									"advanced_retention_type": "YEARLY",
									"retention":               "1200",
								},
							},
							"backup_type":           "COMPLETE",
							"replication_region_id": "cn-beijing",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P2D",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
							"replication_region_id": "cn-shanghai",
						},
						{
							"rule_type":             "REPLICATION",
							"retention":             "8",
							"replication_region_id": "cn-beijing",
							"archive_days":          "50",
							"backup_type":           "COMPLETE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元-修改",
						"policy_name":        name + "_update",
						"rules.#":            "3",
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

// Test Hbr Policy. <<< Resource test cases, automatically generated.
