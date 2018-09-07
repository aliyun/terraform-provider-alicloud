package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ots_instance", &resource.Sweeper{
		Name: "alicloud_ots_instance",
		F:    testSweepOtsInstances,
	})
}

func testSweepOtsInstances(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var insts []ots.InstanceInfo
	req := ots.CreateListInstanceRequest()
	req.RegionId = conn.RegionId
	req.Method = "GET"
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNum = requests.NewInteger(1)
	for {
		resp, err := conn.otsconn.ListInstance(req)
		if err != nil {
			return fmt.Errorf("Error retrieving OTS Instances: %s", err)
		}
		if resp == nil || len(resp.InstanceInfos.InstanceInfo) < 1 {
			break
		}
		insts = append(insts, resp.InstanceInfos.InstanceInfo...)

		if len(resp.InstanceInfos.InstanceInfo) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNum); err != nil {
			return err
		} else {
			req.PageNum = page
		}
	}
	sweeped := false

	for _, v := range insts {
		name := v.InstanceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping OTS Instance: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting OTS Instance: %s", name)
		req := ots.CreateDeleteInstanceRequest()
		req.InstanceName = name
		if _, err := conn.otsconn.DeleteInstance(req); err != nil {
			log.Printf("[ERROR] Failed to delete OTS Instance (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(3 * time.Minute)
	}
	return nil
}

func TestAccAlicloudOtsInstance_Basic(t *testing.T) {
	var instance ots.InstanceInfo
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ots_instance.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist(
						"alicloud_ots_instance.basic", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.basic",
						"name", "tf-testAccBasic"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.basic",
						"accessed_by", "Any"),
				),
			},
		},
	})

}

func TestAccAlicloudOtsInstance_Tags(t *testing.T) {
	var instance ots.InstanceInfo
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ots_instance.tags",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstanceTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist(
						"alicloud_ots_instance.tags", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"name", "tf-testAccTags"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"instance_type", "HighPerformance"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"tags.Created", "TF"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"tags.For", "acceptance test"),
				),
			},
		},
	})

}

func testAccCheckOtsInstanceExist(n string, instance *ots.InstanceInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeOtsInstance(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = &response
		return nil
	}
}

func testAccCheckOtsInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_instance" {
			continue
		}

		client := testAccProvider.Meta().(*AliyunClient)

		if _, err := client.DescribeOtsInstance(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Ots instance %s still exists.", rs.Primary.ID)
	}

	return nil
}

const testAccOtsInstance = `
variable "name" {
  default = "tf-testAccBasic"
}
resource "alicloud_ots_instance" "basic" {
  name = "${var.name}"
  description = "${var.name}"
}
`

const testAccOtsInstanceTags = `
variable "name" {
  default = "tf-testAccTags"
}
resource "alicloud_ots_instance" "tags" {
  name = "${var.name}"
  description = "${var.name}"
  accessed_by = "Vpc"
  tags {
	Created = "TF"
	For = "acceptance test"
  }
}
`
