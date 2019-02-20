package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"log"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_accessrule", &resource.Sweeper{
		Name: "alicloud_nas_accessrule",
		F:    testSweepNasAR,
	})
}

func testSweepNasAR(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
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
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckARDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasArConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckARExists("alicloud_nas_accessrule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "sourcecidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "rwaccess_type", "RDWR"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "useraccess_type", "no_squash"),
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
		CheckDestroy: testAccCheckARDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasArConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckARExists("alicloud_nas_accessrule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "sourcecidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "rwaccess_type", "RDWR"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "useraccess_type", "no_squash"),
				),
			},
			resource.TestStep{
				Config: testAccNasArConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckARExists("alicloud_nas_accessrule.foo", &ar),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "sourcecidr_ip", "172.168.1.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "rwaccess_type", "RDONLY"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessrule.foo", "useraccess_type", "root_squash"),
				),
			},
		},
	})
}


func testAccCheckARExists(n string, nas *nas.AccessRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NAS ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		split := strings.Split(rs.Primary.ID, ":")
		instance, err := nasService.DescribeAccessRules(split[0])

		if err != nil {
			return err
		}

		*nas = instance
		return nil
	}
}

func testAccCheckARDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_accessrule" {
			continue
		}

		// Try to find the NAS
		split := strings.Split(rs.Primary.ID, ":")
		instance, err := nasService.DescribeAccessRules(split[0])

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.AccessRuleId != "" {
			return fmt.Errorf("NAS %s still exist", instance.AccessRuleId)
		}
	}

	return nil
}

const testAccNasArConfig = `
resource "alicloud_nas_accessgroup" "foo" {
		accessgroup_name = "test_wang"
		accessgroup_type = "Classic"
		description = "test_wang"
}
resource "alicloud_nas_accessrule" "foo" {
		accessgroup_name = "${alicloud_nas_accessgroup.foo.id}"
		sourcecidr_ip = "168.1.1.0/16"
		rwaccess_type = "RDWR"
		useraccess_type = "no_squash"
}
`

const testAccNasArConfigUpdate = `
resource "alicloud_nas_accessgroup" "foo" {
		accessgroup_name = "test_wang"
		accessgroup_type = "Classic"
		description = "test_wang"
}
resource "alicloud_nas_accessrule" "foo" {
		accessgroup_name = "${alicloud_nas_accessgroup.foo.id}"
		sourcecidr_ip = "172.168.1.0/16"
		rwaccess_type = "RDONLY"
		useraccess_type = "root_squash"
}
`
