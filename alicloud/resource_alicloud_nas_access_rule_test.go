package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudNasAccessRule_update(t *testing.T) {
	var v nas.AccessRule
	rand := acctest.RandIntRange(10000, 999999)
	resourceID := "alicloud_nas_access_rule.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessRuleVpcConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":    "168.1.1.0/16",
						"rw_access_type":    "RDWR",
						"user_access_type":  "no_squash",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"priority":          "2",
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasAccessRuleConfigUpdateIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":    "172.168.1.0/16",
						"rw_access_type":    "RDWR",
						"user_access_type":  "root_squash",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"priority":          "2",
					}),
				),
			},
			{
				Config: testAccNasAccessRuleConfigUpdateuser_type(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":    "172.168.1.0/16",
						"rw_access_type":    "RDWR",
						"user_access_type":  "all_squash",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"priority":          "2",
					}),
				),
			},
			{
				Config: testAccNasAccessRuleConfigUpdatepriority(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":    "172.168.1.0/16",
						"rw_access_type":    "RDWR",
						"user_access_type":  "all_squash",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"priority":          "10",
					}),
				),
			},
			{
				Config: testAccNasAccessRuleConfigUpdaterw_type(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":    "172.168.1.0/16",
						"rw_access_type":    "RDONLY",
						"user_access_type":  "root_squash",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"priority":          "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudNasAccessRule_Multi(t *testing.T) {
	var v nas.AccessRule
	rand := acctest.RandIntRange(10000, 999999)
	resourceID := "alicloud_nas_access_rule.default.4"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessRuleMulti(rand, acctest.RandIntRange(5, 20)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip":    "168.1.1.0/16",
						"rw_access_type":    "RDWR",
						"user_access_type":  "no_squash",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"priority":          "2",
					}),
				),
			},
		},
	})
}

func testAccCheckAccessRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_access_rule" {
			continue
		}

		// Try to find the NAS
		instance, err := nasService.DescribeNasAccessRule(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if instance.AccessRuleId != "" {
			return WrapError(fmt.Errorf("NAS %s still exist", instance.AccessRuleId))
		}
	}

	return nil
}

func testAccNasAccessRuleVpcConfig(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
        }
        resource "alicloud_nas_access_rule" "default" {
                access_group_name = "${alicloud_nas_access_group.default.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }`, rand)
}

func testAccNasAccessRuleConfigUpdateIp(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_rule" "default" {
		access_group_name = "${alicloud_nas_access_group.default.id}"
                source_cidr_ip = "172.168.1.0/16"
		rw_access_type = "RDWR"
                user_access_type = "root_squash"
		priority = 2
 
	}`, rand)
}

func testAccNasAccessRuleConfigUpdaterw_type(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
        }
        resource "alicloud_nas_access_rule" "default" {
                access_group_name = "${alicloud_nas_access_group.default.id}"
                source_cidr_ip = "172.168.1.0/16"
                rw_access_type = "RDONLY"
                user_access_type = "root_squash"
                priority = 2
 
        }`, rand)
}

func testAccNasAccessRuleConfigUpdateuser_type(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
        }
        resource "alicloud_nas_access_rule" "default" {
                access_group_name = "${alicloud_nas_access_group.default.id}"
                source_cidr_ip = "172.168.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "all_squash"
                priority = 2
 
        }`, rand)
}

func testAccNasAccessRuleConfigUpdatepriority(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
        }
        resource "alicloud_nas_access_rule" "default" {
                access_group_name = "${alicloud_nas_access_group.default.id}"
                source_cidr_ip = "172.168.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "all_squash"
                priority = 10
 
        }`, rand)
}

func testAccNasAccessRuleMulti(rand, count int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
        }
        resource "alicloud_nas_access_rule" "default" {
				count = %d
                access_group_name = "${alicloud_nas_access_group.default.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }`, rand, count)
}
