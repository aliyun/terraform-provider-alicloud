package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudHbrPolicyBinding_basic6295(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6295)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6295)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":    "UDM_ECS",
					"policy_id":      "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id": "${alicloud_instance.defaultrdRDjb.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "UDM_ECS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example (update)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example (update)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "UDM_ECS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_instance.defaultrdRDjb.id}",
					"policy_binding_description": "policy binding example",
					"advanced_options": []map[string]interface{}{
						{
							"udm_detail": []map[string]interface{}{
								{
									"disk_id_list": []string{
										"d-****************zxcv"},
									"destination_kms_key_id": "snxs-******-******-llam",
									"exclude_disk_id_list": []string{
										"d-****************mopl", "d-****************aqlp"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "UDM_ECS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
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

var AliCloudHbrPolicyBindingMap6295 = map[string]string{
	"source_type":    CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
	"policy_id":      CHECKSET,
}

func AliCloudHbrPolicyBindingBasicDependence6295(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
  image_id          = data.alicloud_images.default.images.0.id
  system_disk_category = "cloud_efficiency"
  cpu_core_count                    = 4
  minimum_eni_ipv6_address_quantity = 2
}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  name = var.name
  ipv6_cidr_block_mask = "22"
}

resource "alicloud_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_instance" "defaultrdRDjb" {
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  ipv6_address_count = 1
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_name = var.name
  vswitch_id = "${alicloud_vswitch.vsw.id}"
  internet_max_bandwidth_out = 10
  security_groups = "${alicloud_security_group.group.*.id}"
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
  }
}


`, name)
}

// Case Alibaba Nas Backup 6226
func TestAccAliCloudHbrPolicyBinding_basic6226(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6226)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6226)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":    "NAS",
					"policy_id":      "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id": "${data.alicloud_nas_file_systems.default.systems.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "NAS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example (update)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example (update)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "/backup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "/backup",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "NAS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${data.alicloud_nas_file_systems.default.systems.0.id}",
					"policy_binding_description": "policy binding example",
					"source":                     "/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "NAS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
						"source":                     "/",
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

var AliCloudHbrPolicyBindingMap6226 = map[string]string{
	"source_type":    CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
	"policy_id":      CHECKSET,
}

func AliCloudHbrPolicyBindingBasicDependence6226(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaultTDOTE0" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    keep_latest_snapshots = 1
    vault_id     = alicloud_hbr_vault.defaultTDOTE0.id
  }
}

data "alicloud_nas_file_systems" "default" {
  protocol_type       = "NFS"
}
`, name)
}

// Case OSS Backup 6221
func TestAccAliCloudHbrPolicyBinding_basic6221(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6221)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6221)
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
					"source_type":    "OSS",
					"policy_id":      "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id": "${alicloud_oss_bucket.defaultKtt2XY.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "OSS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "prefix-example-create/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "prefix-example-create/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example (update)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example (update)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "prefix-example-update/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "prefix-example-update/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "OSS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_oss_bucket.defaultKtt2XY.bucket}",
					"policy_binding_description": "policy binding example",
					"source":                     "prefix-example-create/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "OSS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
						"source":                     "prefix-example-create/",
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

var AliCloudHbrPolicyBindingMap6221 = map[string]string{
	"source_type":    CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
	"policy_id":      CHECKSET,
}

func AliCloudHbrPolicyBindingBasicDependence6221(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaultyk84Hc" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    vault_id     = alicloud_hbr_vault.defaultyk84Hc.id
  }
  policy_description = "policy example"
}

resource "alicloud_oss_bucket" "defaultKtt2XY" {
  storage_class = "Standard"
  bucket        = var.name
}


`, name)
}

// Case ECS File Backup 6219
func TestAccAliCloudHbrPolicyBinding_basic6219(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6219)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6219)
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
					"source_type":    "ECS_FILE",
					"policy_id":      "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id": "i-******************",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "ECS_FILE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude": "[\\\"*.pdf\\\",\\\"*.docx\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude": "[\"*.pdf\",\"*.docx\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"include": "[\\\"*.sh\\\",\\\"*.xml\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include": "[\"*.sh\",\"*.xml\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "/root",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "/root",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"speed_limit": "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"speed_limit": "0:24:1024",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example (update)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example (update)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude": "[\\\"*.pdf\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude": "[\"*.pdf\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"include": "[\\\"*.sh\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include": "[\"*.sh\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "/opt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "/opt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"speed_limit": "0:24:2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"speed_limit": "0:24:2048",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "ECS_FILE",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "i-******************",
					"policy_binding_description": "policy binding example",
					"exclude":                    "[\\\"*.pdf\\\",\\\"*.docx\\\"]",
					"include":                    "[\\\"*.sh\\\",\\\"*.xml\\\"]",
					"source":                     "/root",
					"speed_limit":                "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "ECS_FILE",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             "i-******************",
						"policy_binding_description": "policy binding example",
						"exclude":                    "[\"*.pdf\",\"*.docx\"]",
						"include":                    "[\"*.sh\",\"*.xml\"]",
						"source":                     "/root",
						"speed_limit":                "0:24:1024",
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

var AliCloudHbrPolicyBindingMap6219 = map[string]string{
	"source_type":    CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
	"policy_id":      CHECKSET,
}

func AliCloudHbrPolicyBindingBasicDependence6219(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaultQNFISO" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    vault_id     = alicloud_hbr_vault.defaultQNFISO.id
  }
  policy_description = "policy example"
}


`, name)
}

// Case File Backup 6220
func TestAccAliCloudHbrPolicyBinding_basic6220(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6220)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6220)
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
					"source_type":    "File",
					"policy_id":      "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id": "c-******************",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "File",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude": "[\\\"*.pdf\\\",\\\"*.docx\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude": "[\"*.pdf\",\"*.docx\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"include": "[\\\"*.sh\\\",\\\"*.xml\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include": "[\"*.sh\",\"*.xml\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "/root",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "/root",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"speed_limit": "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"speed_limit": "0:24:1024",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_binding_description": "policy binding example (update)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_binding_description": "policy binding example (update)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude": "[\\\"*.pdf\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude": "[\"*.pdf\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"include": "[\\\"*.sh\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include": "[\"*.sh\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "/opt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "/opt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"speed_limit": "0:24:2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"speed_limit": "0:24:2048",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "File",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "c-******************",
					"policy_binding_description": "policy binding example",
					"exclude":                    "[\\\"*.pdf\\\",\\\"*.docx\\\"]",
					"include":                    "[\\\"*.sh\\\",\\\"*.xml\\\"]",
					"source":                     "/root",
					"speed_limit":                "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "File",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             "c-******************",
						"policy_binding_description": "policy binding example",
						"exclude":                    "[\"*.pdf\",\"*.docx\"]",
						"include":                    "[\"*.sh\",\"*.xml\"]",
						"source":                     "/root",
						"speed_limit":                "0:24:1024",
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

var AliCloudHbrPolicyBindingMap6220 = map[string]string{
	"source_type":    CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
	"policy_id":      CHECKSET,
}

func AliCloudHbrPolicyBindingBasicDependence6220(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaultNcs9DB" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|1631685600|P1D"
    retention    = "7"
    archive_days = "0"
    vault_id     = alicloud_hbr_vault.defaultNcs9DB.id
  }
  policy_description = "policy example"
}


`, name)
}

// Case ECS Instance Backup 6295  twin
func TestAccAliCloudHbrPolicyBinding_basic6295_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6295)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6295)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "UDM_ECS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_instance.defaultrdRDjb.id}",
					"policy_binding_description": "policy binding example (update)",
					"advanced_options": []map[string]interface{}{
						{
							"udm_detail": []map[string]interface{}{
								{
									"disk_id_list": []string{
										"d-****************lpsa", "d-****************09", "d-****************qpla"},
									"destination_kms_key_id": "sdak-******-******-qozp",
									"exclude_disk_id_list": []string{
										"d-****************qqaa", "d-****************qqxv", "d-****************qblz"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "UDM_ECS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example (update)",
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

// Case Alibaba Nas Backup 6226  twin
func TestAccAliCloudHbrPolicyBinding_basic6226_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6226)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6226)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "NAS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${data.alicloud_nas_file_systems.default.systems.0.id}",
					"policy_binding_description": "policy binding example (update)",
					"source":                     "/backup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "NAS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example (update)",
						"source":                     "/backup",
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

// Case OSS Backup 6221  twin
func TestAccAliCloudHbrPolicyBinding_basic6221_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6221)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6221)
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
					"source_type":                "OSS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_oss_bucket.defaultKtt2XY.bucket}",
					"policy_binding_description": "policy binding example (update)",
					"source":                     "prefix-example-update/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "OSS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example (update)",
						"source":                     "prefix-example-update/",
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

// Case ECS File Backup 6219  twin
func TestAccAliCloudHbrPolicyBinding_basic6219_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6219)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6219)
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
					"source_type":                "ECS_FILE",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "i-******************",
					"policy_binding_description": "policy binding example (update)",
					"exclude":                    "[\\\"*.pdf\\\"]",
					"include":                    "[\\\"*.sh\\\"]",
					"source":                     "/opt",
					"speed_limit":                "0:24:2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "ECS_FILE",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             "i-******************",
						"policy_binding_description": "policy binding example (update)",
						"exclude":                    "[\"*.pdf\"]",
						"include":                    "[\"*.sh\"]",
						"source":                     "/opt",
						"speed_limit":                "0:24:2048",
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

// Case File Backup 6220  twin
func TestAccAliCloudHbrPolicyBinding_basic6220_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6220)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6220)
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
					"source_type":                "File",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "c-******************",
					"policy_binding_description": "policy binding example (update)",
					"exclude":                    "[\\\"*.pdf\\\"]",
					"include":                    "[\\\"*.sh\\\"]",
					"source":                     "/opt",
					"speed_limit":                "0:24:2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "File",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             "c-******************",
						"policy_binding_description": "policy binding example (update)",
						"exclude":                    "[\"*.pdf\"]",
						"include":                    "[\"*.sh\"]",
						"source":                     "/opt",
						"speed_limit":                "0:24:2048",
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

// Case ECS Instance Backup 6295   raw
func TestAccAliCloudHbrPolicyBinding_basic6295_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6295)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6295)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "UDM_ECS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_instance.defaultrdRDjb.id}",
					"policy_binding_description": "policy binding example",
					"advanced_options": []map[string]interface{}{
						{
							"udm_detail": []map[string]interface{}{
								{
									"disk_id_list": []string{
										"d-****************zxcv"},
									"destination_kms_key_id": "snxs-******-******-llam",
									"exclude_disk_id_list": []string{
										"d-****************mopl", "d-****************aqlp"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "UDM_ECS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled":                   "true",
					"policy_binding_description": "policy binding example (update)",
					"advanced_options": []map[string]interface{}{
						{
							"udm_detail": []map[string]interface{}{
								{
									"disk_id_list": []string{
										"d-****************lpsa", "d-****************09", "d-****************qpla"},
									"destination_kms_key_id": "sdak-******-******-qozp",
									"exclude_disk_id_list": []string{
										"d-****************qqaa", "d-****************qqxv", "d-****************qblz"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled":                   "true",
						"policy_binding_description": "policy binding example (update)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
					//"advanced_options": []map[string]interface{}{
					//	{
					//		"udm_detail": []map[string]interface{}{
					//			{
					//				"disk_id_list":         []string{},
					//				"exclude_disk_id_list": []string{},
					//			},
					//		},
					//	},
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
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

// Case Alibaba Nas Backup 6226   raw
func TestAccAliCloudHbrPolicyBinding_basic6226_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6226)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6226)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":                "NAS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${data.alicloud_nas_file_systems.default.systems.0.id}",
					"policy_binding_description": "policy binding example",
					"source":                     "/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "NAS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
						"source":                     "/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled":                   "true",
					"policy_binding_description": "policy binding example (update)",
					"source":                     "/backup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled":                   "true",
						"policy_binding_description": "policy binding example (update)",
						"source":                     "/backup",
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

// Case OSS Backup 6221   raw
func TestAccAliCloudHbrPolicyBinding_basic6221_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6221)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6221)
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
					"source_type":                "OSS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_oss_bucket.defaultKtt2XY.bucket}",
					"policy_binding_description": "policy binding example",
					"source":                     "prefix-example-create/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "OSS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
						"source":                     "prefix-example-create/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled":                   "true",
					"policy_binding_description": "policy binding example (update)",
					"source":                     "prefix-example-update/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled":                   "true",
						"policy_binding_description": "policy binding example (update)",
						"source":                     "prefix-example-update/",
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

// Case ECS File Backup 6219   raw
func TestAccAliCloudHbrPolicyBinding_basic6219_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6219)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6219)
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
					"source_type":                "ECS_FILE",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "i-******************",
					"policy_binding_description": "policy binding example",
					"exclude":                    "[\\\"*.pdf\\\",\\\"*.docx\\\"]",
					"include":                    "[\\\"*.sh\\\",\\\"*.xml\\\"]",
					"source":                     "/root",
					"speed_limit":                "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "ECS_FILE",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             "i-******************",
						"policy_binding_description": "policy binding example",
						"exclude":                    "[\"*.pdf\",\"*.docx\"]",
						"include":                    "[\"*.sh\",\"*.xml\"]",
						"source":                     "/root",
						"speed_limit":                "0:24:1024",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled":                   "true",
					"policy_binding_description": "policy binding example (update)",
					"exclude":                    "[\\\"*.pdf\\\"]",
					"include":                    "[\\\"*.sh\\\"]",
					"source":                     "/opt",
					"speed_limit":                "0:24:2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled":                   "true",
						"policy_binding_description": "policy binding example (update)",
						"exclude":                    "[\"*.pdf\"]",
						"include":                    "[\"*.sh\"]",
						"source":                     "/opt",
						"speed_limit":                "0:24:2048",
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

// Case File Backup 6220   raw
func TestAccAliCloudHbrPolicyBinding_basic6220_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap6220)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicybinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence6220)
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
					"source_type":                "File",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "c-******************",
					"policy_binding_description": "policy binding example",
					"exclude":                    "[\\\"*.pdf\\\",\\\"*.docx\\\"]",
					"include":                    "[\\\"*.sh\\\",\\\"*.xml\\\"]",
					"source":                     "/root",
					"speed_limit":                "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "File",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             "c-******************",
						"policy_binding_description": "policy binding example",
						"exclude":                    "[\"*.pdf\",\"*.docx\"]",
						"include":                    "[\"*.sh\",\"*.xml\"]",
						"source":                     "/root",
						"speed_limit":                "0:24:1024",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled":                   "true",
					"policy_binding_description": "policy binding example (update)",
					"exclude":                    "[\\\"*.pdf\\\"]",
					"include":                    "[\\\"*.sh\\\"]",
					"source":                     "/opt",
					"speed_limit":                "0:24:2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled":                   "true",
						"policy_binding_description": "policy binding example (update)",
						"exclude":                    "[\"*.pdf\"]",
						"include":                    "[\"*.sh\"]",
						"source":                     "/opt",
						"speed_limit":                "0:24:2048",
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

// Case OSS Backup跨账号 7232
func TestAccAliCloudHbrPolicyBinding_basic7232(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_policy_binding.default.19"
	ra := resourceAttrInit(resourceId, AliCloudHbrPolicyBindingMap7232)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrPolicyBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrpolicyossbind%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudHbrPolicyBindingBasicDependence7232)
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
					"count":                      "20",
					"source_type":                "OSS",
					"disabled":                   "false",
					"policy_id":                  "${alicloud_hbr_policy.defaultoqWvHQ.id}",
					"data_source_id":             "${alicloud_oss_bucket.defaultKtt2XY[count.index].bucket}",
					"policy_binding_description": "policy binding example",
					"source":                     "prefix-example-create/",
					"cross_account_user_id":      "1",
					"cross_account_role_name":    "mock",
					"cross_account_type":         "CROSS_ACCOUNT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":                "OSS",
						"disabled":                   "false",
						"policy_id":                  CHECKSET,
						"data_source_id":             CHECKSET,
						"policy_binding_description": "policy binding example",
						"source":                     "prefix-example-create/",
						"cross_account_user_id":      "1",
						"cross_account_role_name":    "mock",
						"cross_account_type":         "CROSS_ACCOUNT",
					}),
				),
			},
		},
	})
}

var AliCloudHbrPolicyBindingMap7232 = map[string]string{
	"create_time":    CHECKSET,
	"source_type":    CHECKSET,
	"policy_id":      CHECKSET,
	"data_source_id": CHECKSET,
}

func AliCloudHbrPolicyBindingBasicDependence7232(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_hbr_vault" "defaultyk84Hc" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_description = "created by zhenyuan"
  policy_name        = format("%%s1", var.name)
  rules {
    rule_type             = "BACKUP"
    backup_type           = "COMPLETE"
    schedule              = "I|1631685600|P1D"
    retention             = "7"
    vault_id              = alicloud_hbr_vault.defaultyk84Hc.id
    keep_latest_snapshots = "1"
    archive_days          = "0"
  }
}

resource "alicloud_oss_bucket" "defaultKtt2XY" {
  count         = 20
  storage_class = "Standard"
  bucket        = format("%%s2%%s", var.name, count.index)
}


`, name)
}
