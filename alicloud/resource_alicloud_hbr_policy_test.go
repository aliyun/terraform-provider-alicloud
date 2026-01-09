package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Hbr Policy. >>> Resource test cases, automatically generated.
// Case Policy测试用例 5320
func TestAccAliCloudHbrPolicy_basic5320(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap5320)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence5320)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
							"rule_type":             "TRANSITION",
							"retention":             "120",
							"archive_days":          "30",
							"replication_region_id": "cn-shanghai",
							"keep_latest_snapshots": "1",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "1",
							"archive_days":          "30",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
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
							"rule_type":             "TRANSITION",
							"retention":             "240",
							"archive_days":          "85",
							"backup_type":           "COMPLETE",
							"replication_region_id": "cn-beijing",
							"keep_latest_snapshots": "1",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P2D",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "85",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元-修改",
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

var AlicloudHbrPolicyMap5320 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence5320(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaulth4dKAG" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "ault-example-1767865533"
  vault_storage_class = "STANDARD"
}


`, name)
}

// Case Policy 6287
func TestAccAliCloudHbrPolicy_basic6287(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap6287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence6287)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
						{
							"rule_type":    "BACKUP",
							"backup_type":  "INCREMENTAL",
							"schedule":     "I|1631685600|P1M",
							"vault_id":     "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"archive_days": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name,
						"rules.#":            "4",
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
							"rule_type":             "TRANSITION",
							"retention":             "120",
							"archive_days":          "0",
							"replication_region_id": "cn-chengdu",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"replication_region_id": "cn-chengdu",
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
							"rule_type":    "TRANSITION",
							"retention":    "145",
							"archive_days": "0",
						},
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
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

var AlicloudHbrPolicyMap6287 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence6287(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaulth4dKAG" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "vault-example-1767865534"
  vault_storage_class = "STANDARD"
}

resource "alicloud_hbr_vault" "defaultL7kwwD" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "vault-example-1767865535"
  vault_storage_class = "STANDARD"
}


`, name)
}

// Case Policy 6963
func TestAccAliCloudHbrPolicy_basic6963(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap6963)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence6963)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
								{
									"advanced_retention_type": "DAILY",
									"retention":               "170",
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
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name,
						"rules.#":            "2",
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
							"retention":    "2",
							"archive_days": "85",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "DAILY",
									"retention":               "920",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "1000",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|PT1H",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
							"retention":             "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy update",
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

var AlicloudHbrPolicyMap6963 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence6963(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaulth4dKAG" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "vault-example-1767865535"
  vault_storage_class = "STANDARD"
}


`, name)
}

// Case PolicyV1.2.0 7188
func TestAccAliCloudHbrPolicy_basic7188(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap7188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence7188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
							"rule_type":             "TRANSITION",
							"retention":             "120",
							"archive_days":          "30",
							"replication_region_id": "cn-shanghai",
							"keep_latest_snapshots": "1",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"keep_latest_snapshots": "1",
							"archive_days":          "30",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
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
							"rule_type":             "TRANSITION",
							"retention":             "240",
							"archive_days":          "85",
							"backup_type":           "COMPLETE",
							"replication_region_id": "cn-beijing",
							"keep_latest_snapshots": "1",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P2D",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "85",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "镇元-修改",
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

var AlicloudHbrPolicyMap7188 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence7188(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaulth4dKAG" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "ault-example-1767865536"
  vault_storage_class = "STANDARD"
}


`, name)
}

// Case PolicyV1.2.0 7189
func TestAccAliCloudHbrPolicy_basic7189(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap7189)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence7189)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
								{
									"advanced_retention_type": "DAILY",
									"retention":               "170",
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
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name,
						"rules.#":            "2",
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
							"retention":    "2",
							"archive_days": "85",
							"retention_rules": []map[string]interface{}{
								{
									"advanced_retention_type": "DAILY",
									"retention":               "920",
								},
								{
									"advanced_retention_type": "MONTHLY",
									"retention":               "1000",
								},
							},
							"backup_type": "COMPLETE",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|PT1H",
							"vault_id":              "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"keep_latest_snapshots": "1",
							"archive_days":          "0",
							"retention":             "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy update",
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

var AlicloudHbrPolicyMap7189 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence7189(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaulth4dKAG" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "vault-example-1767865536"
  vault_storage_class = "STANDARD"
}


`, name)
}

// Case ECS整机备份策略 10119
func TestAccAliCloudHbrPolicy_basic10119(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap10119)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence10119)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
						},
					},
					"policy_type": "UDM_ECS_ONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name,
						"rules.#":            "2",
						"policy_type":        "UDM_ECS_ONLY",
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

var AlicloudHbrPolicyMap10119 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence10119(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 通用备份策略 7187
func TestAccAliCloudHbrPolicy_basic7187(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrPolicyMap7187)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrPolicyBasicDependence7187)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
						{
							"rule_type":    "BACKUP",
							"backup_type":  "INCREMENTAL",
							"schedule":     "I|1631685600|P1M",
							"vault_id":     "${alicloud_hbr_vault.defaulth4dKAG.id}",
							"archive_days": "0",
						},
					},
					"policy_type": "STANDARD",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_description": "policy creation",
						"policy_name":        name,
						"rules.#":            "4",
						"policy_type":        "STANDARD",
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
							"rule_type":             "TRANSITION",
							"retention":             "120",
							"archive_days":          "0",
							"replication_region_id": "cn-chengdu",
						},
						{
							"rule_type":             "BACKUP",
							"backup_type":           "COMPLETE",
							"schedule":              "I|1631685600|P1D",
							"archive_days":          "0",
							"vault_id":              "${alicloud_hbr_vault.defaultL7kwwD.id}",
							"replication_region_id": "cn-chengdu",
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
							"rule_type":    "TRANSITION",
							"retention":    "145",
							"archive_days": "0",
						},
						{
							"rule_type":    "BACKUP",
							"backup_type":  "COMPLETE",
							"schedule":     "I|1631685600|P1D",
							"vault_id":     "${alicloud_hbr_vault.defaultL7kwwD.id}",
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

var AlicloudHbrPolicyMap7187 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudHbrPolicyBasicDependence7187(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaulth4dKAG" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "vault-example-1767865538"
  vault_storage_class = "STANDARD"
}

resource "alicloud_hbr_vault" "defaultL7kwwD" {
  vault_type          = "STANDARD"
  encrypt_type        = "HBR_PRIVATE"
  vault_name          = "vault-example-1767865539"
  vault_storage_class = "STANDARD"
}


`, name)
}

// Test Hbr Policy. <<< Resource test cases, automatically generated.
