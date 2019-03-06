package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_access_rule", &resource.Sweeper{
		Name: "alicloud_nas_access_rule",
		F:    testSweepNasRule,
	})
}

func testSweepNasRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var ar []nas.AccessRule
	req := nas.CreateDescribeAccessRulesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessRules(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving filesystem: %s", err)
		}
		resp, _ := raw.(*nas.DescribeAccessRulesResponse)
		if resp == nil || len(resp.AccessRules.AccessRule) < 1 {
			break
		}
		ar = append(ar, resp.AccessRules.AccessRule...)

		if len(resp.AccessRules.AccessRule) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, fs := range ar {

		id := fs.AccessRuleId
		SourceCidrIp := fs.SourceCidrIp
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(SourceCidrIp), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping AccessRule: %s (%s)", SourceCidrIp, id)
			continue
		}
		log.Printf("[INFO] Deleting AccessRule: %s (%s)", SourceCidrIp, id)
		req := nas.CreateDeleteAccessRuleRequest()
		req.AccessGroupName = id
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteAccessRule(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete AccessRule (%s (%s)): %s", SourceCidrIp, id, err)
		}
	}
	return nil
}

func TestAccAlicloudNas_AccessRule_basic(t *testing.T) {
	var ar nas.AccessRule
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAccessRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessRuleExists("alicloud_nas_access_rule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "rw_access_type", "RDWR"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "user_access_type", "no_squash"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "access_group_name", "tf-testAccNasConfigName"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "priority", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudNas_AccessRule_update(t *testing.T) {
	var ar nas.AccessRule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAccessRuleConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessRuleExists("alicloud_nas_access_rule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "access_group_name", "tf-testAccNasConfigName"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "rw_access_type", "RDWR"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "user_access_type", "no_squash"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "priority", "2"),
				),
			},
			resource.TestStep{
				Config: testAccNasAccessRuleConfigUpdateIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccessRuleExists("alicloud_nas_access_rule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "source_cidr_ip", "172.168.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_rule.foo", "access_group_name", "tf-testAccNasConfigName"),
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

func testAccCheckAccessRuleExists(n string, nas *nas.AccessRule) resource.TestCheckFunc {
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

const testAccNasAccessRuleConfig = `
resource "alicloud_nas_access_group" "foo" {
		name = "tf-testAccNasConfigName"
		type = "Classic"
		description = "tf-testAccNasConfig"
}
resource "alicloud_nas_access_rule" "foo" {
		access_group_name = "${alicloud_nas_access_group.foo.id}"
		source_cidr_ip = "168.1.1.0/16"
		rw_access_type = "RDWR"
		user_access_type = "no_squash"
		priority = 2
}
`

const testAccNasAccessRuleConfigUpdateIp = `
resource "alicloud_nas_access_group" "foo" {
                name = "tf-testAccNasConfigName"
                type = "Classic"
                description = "tf-testAccNasConfig"
}
resource "alicloud_nas_access_rule" "foo" {
		access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "172.168.1.0/16"
		rw_access_type = "RDONLY"
                user_access_type = "root_squash"
		priority = 2
 
}
`
