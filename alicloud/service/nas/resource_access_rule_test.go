package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNasAccessRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNasAccessRule0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasAccessRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNasAccessRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example.access_group_name}",
					"source_cidr_ip":    "168.1.1.0/16",
					"priority":          "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"source_cidr_ip":    "168.1.1.0/16",
						"priority":          "1",
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
					"source_cidr_ip": "172.168.1.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip": "172.168.1.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rw_access_type": "RDONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_access_type": "RDONLY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type": "root_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type": "root_squash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_cidr_ip": "168.1.1.0/16",
					"priority":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip": "168.1.1.0/16",
						"priority":       "1",
					}),
				),
			},
		},
	})
}

var AliCloudNasAccessRule0 = map[string]string{
	"access_rule_id": acctest.CHECKSET,
	"region_id":      acctest.CHECKSET,
}

func AliCloudNasAccessRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
	}
`, name)
}

// lintignore: R001
func TestAccAliCloudNasAccessRule_basic5737(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNasAccessRuleMap5737)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccessrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNasAccessRuleBasicDependence5737)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.NasExtremeTypeRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_type":    "extreme",
					"access_group_name":   "${alicloud_nas_access_group.defaultf20OX7.access_group_name}",
					"ipv6_source_cidr_ip": "2001:250:6000::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_type":    "extreme",
						"access_group_name":   acctest.CHECKSET,
						"ipv6_source_cidr_ip": "2001:250:6000::1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type": "no_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type": "no_squash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rw_access_type": "RDONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_access_type": "RDONLY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_cidr_ip":      "",
					"ipv6_source_cidr_ip": "::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":      "",
						"ipv6_source_cidr_ip": "::1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type": "root_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type": "root_squash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rw_access_type": "RDWR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_access_type": "RDWR",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_source_cidr_ip": "2001:250:6000::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_source_cidr_ip": "2001:250:6000::1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type":    "no_squash",
					"priority":            "1",
					"file_system_type":    "extreme",
					"access_group_name":   "${alicloud_nas_access_group.defaultf20OX7.access_group_name}",
					"rw_access_type":      "RDONLY",
					"ipv6_source_cidr_ip": "::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type":    "no_squash",
						"priority":            "1",
						"file_system_type":    "extreme",
						"access_group_name":   acctest.CHECKSET,
						"rw_access_type":      "RDONLY",
						"ipv6_source_cidr_ip": "::1",
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

var AliCloudNasAccessRuleMap5737 = map[string]string{
	"access_rule_id": acctest.CHECKSET,
	"region_id":      acctest.CHECKSET,
}

func AliCloudNasAccessRuleBasicDependence5737(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_nas_access_group" "defaultf20OX7" {
  access_group_type = "Vpc"
  description       = "ExtremeAccessGroup"
  access_group_name = var.name

  file_system_type = "extreme"
}

`, name)
}

// Case 通用型 Vpc 权限组规则测试 5736
func TestAccAliCloudNasAccessRule_basic5736(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNasAccessRuleMap5736)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccessrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNasAccessRuleBasicDependence5736)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_type":  "standard",
					"access_group_name": "${alicloud_nas_access_group.defaultCWVMZb.access_group_name}",
					"source_cidr_ip":    "1.1.1.1/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_type":  "standard",
						"access_group_name": acctest.CHECKSET,
						"source_cidr_ip":    "1.1.1.1/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type": "no_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type": "no_squash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_cidr_ip": "1.1.1.1/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip": "1.1.1.1/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type":  "root_squash",
					"priority":          "2",
					"file_system_type":  "standard",
					"source_cidr_ip":    "1.1.1.2/0",
					"access_group_name": "${alicloud_nas_access_group.defaultCWVMZb.access_group_name}",
					"rw_access_type":    "RDWR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type":  "root_squash",
						"priority":          "2",
						"file_system_type":  "standard",
						"source_cidr_ip":    "1.1.1.2/0",
						"access_group_name": acctest.CHECKSET,
						"rw_access_type":    "RDWR",
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

var AliCloudNasAccessRuleMap5736 = map[string]string{
	"access_rule_id": acctest.CHECKSET,
	"region_id":      acctest.CHECKSET,
}

func AliCloudNasAccessRuleBasicDependence5736(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_nas_access_group" "defaultCWVMZb" {
  access_group_type = "Vpc"
  description       = "StandardAccessGroup"
  access_group_name = var.name

  file_system_type = "standard"
}

`, name)
}

// Case 极速型 Vpc 权限组规则测试 5737  twin
func TestAccAliCloudNasAccessRule_basic5737_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNasAccessRuleMap5737)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccessrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNasAccessRuleBasicDependence5737)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.NasExtremeTypeRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type":    "no_squash",
					"priority":            "1",
					"file_system_type":    "extreme",
					"access_group_name":   "${alicloud_nas_access_group.defaultf20OX7.access_group_name}",
					"rw_access_type":      "RDONLY",
					"ipv6_source_cidr_ip": "::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type":    "no_squash",
						"priority":            "1",
						"file_system_type":    "extreme",
						"access_group_name":   acctest.CHECKSET,
						"rw_access_type":      "RDONLY",
						"ipv6_source_cidr_ip": "::1",
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

// Case 通用型 Vpc 权限组规则测试 5736  twin
func TestAccAliCloudNasAccessRule_basic5736_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNasAccessRuleMap5736)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccessrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNasAccessRuleBasicDependence5736)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type":  "root_squash",
					"priority":          "2",
					"file_system_type":  "standard",
					"source_cidr_ip":    "1.1.1.2/0",
					"access_group_name": "${alicloud_nas_access_group.defaultCWVMZb.access_group_name}",
					"rw_access_type":    "RDWR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type":  "root_squash",
						"priority":          "2",
						"file_system_type":  "standard",
						"source_cidr_ip":    "1.1.1.2/0",
						"access_group_name": acctest.CHECKSET,
						"rw_access_type":    "RDWR",
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

// Test Nas AccessRule. <<< Resource test cases, automatically generated.
