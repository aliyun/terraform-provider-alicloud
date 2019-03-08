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
	resource.AddTestSweepers("alicloud_nas_access_group", &resource.Sweeper{
		Name: "alicloud_nas_access_group",
		F:    testSweepNasAccessGroup,
	})
}

func testSweepNasAccessGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var ag []nas.AccessGroup
	req := nas.CreateDescribeAccessGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving filesystem: %s", err)
		}
		resp, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if resp == nil || len(resp.AccessGroups.AccessGroup) < 1 {
			break
		}
		ag = append(ag, resp.AccessGroups.AccessGroup...)

		if len(resp.AccessGroups.AccessGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, fs := range ag {

		id := fs.AccessGroupName
		AccessGroupType := fs.AccessGroupType
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(AccessGroupType), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping AccessGroup: %s (%s)", AccessGroupType, id)
			continue
		}
		log.Printf("[INFO] Deleting AccessGroup: %s (%s)", AccessGroupType, id)
		req := nas.CreateDeleteAccessGroupRequest()
		req.AccessGroupName = id
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteAccessGroup(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete AccessGroup (%s (%s)): %s", AccessGroupType, id, err)
		}
	}
	return nil
}

func TestAccAlicloudNas_AccessGroup_basic(t *testing.T) {
	var ag nas.AccessGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAGDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAgConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAGExists("alicloud_nas_access_group.foo", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "type", "Classic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "description", "tf-testAccNasConfigDescription"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "name", "tf-testAccNasConfigName"),
				),
			},
		},
	})

}

func TestAccAlicloudNas_AccessGroup_update(t *testing.T) {
	var ag nas.AccessGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAGDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAgConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAGExists("alicloud_nas_access_group.foo", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "type", "Classic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "description", "tf-testAccNasConfigDescription"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "name", "tf-testAccNasConfigName"),
				),
			},
			resource.TestStep{
				Config: testAccNasAgConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAGExists("alicloud_nas_access_group.foo", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.foo", "description", "tf-testAccNasConfigUpdateDescription"),
				),
			},
		},
	})
}

func TestAccAlicloudNas_AccessGroup_multi(t *testing.T) {
	var ag nas.AccessGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAGDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasAgConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAGExists("alicloud_nas_access_group.bar_1", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.bar_1", "name", "tf-testAccNasConfigClassic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.bar_1", "type", "Classic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.bar_1", "description", "tf-testAccNasConfigDescription-1"),
					testAccCheckAGExists("alicloud_nas_access_group.bar_2", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.bar_2", "name", "tf-testAccNasConfigVpc"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.bar_2", "type", "Vpc"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_access_group.bar_2", "description", "tf-testAccNasConfigDescription-2"),
				),
			},
		},
	})
}

func testAccCheckAGExists(n string, nas *nas.AccessGroup) resource.TestCheckFunc {
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
		instance, err := nasService.DescribeNasAccessGroup(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*nas = instance
		return nil
	}
}

func testAccCheckAGDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_file_system" {
			continue
		}

		// Try to find the NAS
		instance, err := nasService.DescribeNasAccessGroup(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("NAS %s still exist", instance.AccessGroupName))
	}

	return nil
}

const testAccNasAgConfig = `
resource "alicloud_nas_access_group" "foo" {
		name = "tf-testAccNasConfigName"
		type = "Classic"
		description = "tf-testAccNasConfigDescription"
}
`

const testAccNasAgConfigUpdate = `
resource "alicloud_nas_access_group" "foo" {
		name = "tf-testAccNasConfigName"
		type = "Classic"
		description = "tf-testAccNasConfigUpdateDescription"
}
`

const testAccNasAgConfigMulti = `
variable "description" {
  	default = "tf-testAccNasConfigDescription"
}
resource "alicloud_nas_access_group" "bar_1" {
	name = "tf-testAccNasConfigClassic"
	type = "Classic"
	description = "${var.description}-1"
}
resource "alicloud_nas_access_group" "bar_2" {
	name = "tf-testAccNasConfigVpc"
	type = "Vpc"
	description = "${var.description}-2"
}
`
