package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNas_AccessRule_basic(t *testing.T) {
	var ar nas.DescribeAccessRulesAccessRule1
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAccessRuleConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessRuleExists("alicloud_nas_access_rule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "rw_access_type", "RDWR"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "user_access_type", "no_squash"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfigName-%d", rand)),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "priority", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudNas_AccessRule_update(t *testing.T) {
	var ar nas.DescribeAccessRulesAccessRule1
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAccessRuleVpcConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessRuleExists("alicloud_nas_access_rule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfigName-%d", rand)),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "rw_access_type", "RDWR"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "user_access_type", "no_squash"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "priority", "2"),
				),
			},
			resource.TestStep{
				Config: testAccNasAccessRuleConfigUpdateIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessRuleExists("alicloud_nas_access_rule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "source_cidr_ip", "172.168.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfigName-%d", rand)),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "rw_access_type", "RDONLY"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "user_access_type", "root_squash"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "priority", "2"),
				),
			},
		},
	})
}

func testAccCheckAccessRuleExists(n string, nas *nas.DescribeAccessRulesAccessRule1) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No NAS ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		split := strings.Split(rs.Primary.ID, ":")
		instance, err := nasService.DescribeNasAccessRule(split[0])

		if err != nil {
			return WrapError(err)
		}

		*nas = instance
		return nil
	}
}

func testAccCheckAccessRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_access_rule" {
			continue
		}

		// Try to find the NAS
		split := strings.Split(rs.Primary.ID, ":")
		instance, err := nasService.DescribeNasAccessRule(split[0])

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

func testAccNasAccessRuleConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "foo" {
		name = "tf-testAccNasConfigName-%d"
		type = "Classic"
		description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_rule" "foo" {
		access_group_name = "${alicloud_nas_access_group.foo.id}"
		source_cidr_ip = "168.1.1.0/16"
		rw_access_type = "RDWR"
		user_access_type = "no_squash"
		priority = 2
	}`, rand)
}

func testAccNasAccessRuleVpcConfig(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "foo" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }`, rand)
}

func testAccNasAccessRuleConfigUpdateIp(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "foo" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_rule" "foo" {
		access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "172.168.1.0/16"
		rw_access_type = "RDONLY"
                user_access_type = "root_squash"
		priority = 2
 
	}`, rand)
}
