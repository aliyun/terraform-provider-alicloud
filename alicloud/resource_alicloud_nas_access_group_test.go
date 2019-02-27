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
	resource.AddTestSweepers("alicloud_nas_accessgroup", &resource.Sweeper{
		Name: "alicloud_nas_access_group",
		F:    testSweepNasAG,
	})
}

func testSweepNasAG(region string) error {
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
					testAccCheckAGExists("alicloud_nas_accessgroup.foo", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "accessgroup_type", "Classic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "description", "test_wang"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "accessgroup_name", "test_wang"),
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
					testAccCheckAGExists("alicloud_nas_accessgroup.foo", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "accessgroup_type", "Classic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "description", "test_wang"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "accessgroup_name", "test_wang"),
				),
			},
			resource.TestStep{
				Config: testAccNasAgConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAGExists("alicloud_nas_accessgroup.foo", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.foo", "description", "wang_test"),
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
					testAccCheckAGExists("alicloud_nas_accessgroup.bar_1", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.bar_1", "accessgroup_name", "test_wang_classic"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.bar_1", "accessgroup_type", "Classic"),
					testAccCheckAGExists("alicloud_nas_accessgroup.bar_2", &ag),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.bar_2", "accessgroup_name", "test_wang_vpc"),
					resource.TestCheckResourceAttr(
						"alicloud_nas_accessgroup.bar_2", "accessgroup_type", "Vpc"),
				),
			},
		},
	})
}

func testAccCheckAGExists(n string, nas *nas.AccessGroup) resource.TestCheckFunc {
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
		instance, err := nasService.DescribeAccessGroup(rs.Primary.ID)

		if err != nil {
			return err
		}

		*nas = instance
		return nil
	}
}

func testAccCheckAGDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_filesystem" {
			continue
		}

		// Try to find the NAS
		instance, err := nasService.DescribeAccessGroup(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.AccessGroupName != "" {
			return fmt.Errorf("NAS %s still exist", instance.AccessGroupName)
		}
	}

	return nil
}

const testAccNasAgConfig = `
resource "alicloud_nas_accessgroup" "foo" {
		accessgroup_name = "test_wang"
		accessgroup_type = "Classic"
		description = "test_wang"
}
`

const testAccNasAgConfigUpdate = `
resource "alicloud_nas_accessgroup" "foo" {
		accessgroup_name = "test_wang"
		accessgroup_type = "Classic"
		description = "wang_test"
}
`

const testAccNasAgConfigMulti = `
variable "description" {
  	default = "test_wang"
}
resource "alicloud_nas_accessgroup" "bar_1" {
	accessgroup_name = "test_wang_classic"
	accessgroup_type = "Classic"
	description = "${var.description}-1"
}
resource "alicloud_nas_accessgroup" "bar_2" {
	accessgroup_name = "test_wang_vpc"
	accessgroup_type = "Vpc"
	description = "${var.description}-2"
}
`
