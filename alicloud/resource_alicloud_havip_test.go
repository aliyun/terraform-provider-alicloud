package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// At present, only white list users can operate HaVip Resource. So close havip sweeper.
//func init() {
//	resource.AddTestSweepers("alicloud_havip", &resource.Sweeper{
//		Name: "alicloud_havip",
//		F:    testSweepHaVip,
//		// When implemented, these should be removed firstly
//		Dependencies: []string{
//			"alicloud_havip_attachment",
//		},
//	})
//}
//
//func testSweepHaVip(region string) error {
//	client, err := sharedClientForRegion(region)
//	if err != nil {
//		return fmt.Errorf("error getting Alicloud client: %s", err)
//	}
//	conn := client.(*AliyunClient)
//
//	prefixes := []string{
//		"tf-testAcc",
//		"tf_testAcc",
//	}
//
//	var haVips []vpc.HaVip
//	req := vpc.CreateDescribeHaVipsRequest()
//	req.RegionId = conn.RegionId
//	req.PageSize = requests.NewInteger(PageSizeLarge)
//	req.PageNumber = requests.NewInteger(1)
//	for {
//		resp, err := conn.vpcconn.DescribeHaVips(req)
//		if err != nil {
//			return fmt.Errorf("Error retrieving HaVips: %s", err)
//		}
//		if resp == nil || len(resp.HaVips.HaVip) < 1 {
//			break
//		}
//		haVips = append(haVips, resp.HaVips.HaVip...)
//
//		if len(resp.HaVips.HaVip) < PageSizeLarge {
//			break
//		}
//
//		if page, err := getNextpageNumber(req.PageNumber); err != nil {
//			return err
//		} else {
//			req.PageNumber = page
//		}
//	}
//
//	for _, havip := range haVips {
//		description := havip.Description
//		id := havip.HaVipId
//		skip := true
//		for _, prefix := range prefixes {
//			if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
//				skip = false
//				break
//			}
//		}
//		if skip {
//			log.Printf("[INFO] Skipping HaVip: %s (%s)", description, id)
//			continue
//		}
//		log.Printf("[INFO] Deleting HaVip: %s (%s)", description, id)
//		req := vpc.CreateDeleteHaVipRequest()
//		req.HaVipId = id
//		if _, err := conn.vpcconn.DeleteHaVip(req); err != nil {
//			log.Printf("[ERROR] Failed to delete HaVip (%s (%s)): %s", description, id, err)
//		}
//	}
//	return nil
//}

// At present, only white list users can operate HaVip Resource.
func SkipTestAccAlicloudHaVip_basic(t *testing.T) {
	var havip vpc.HaVip
	resourceId := "alicloud_havip.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckHaVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipExists("alicloud_havip.foo", havip),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip.foo", "vswitch_id"),
					resource.TestCheckResourceAttr(
						"alicloud_havip.foo", "description", "tf_testAcc_havip"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func SkipTestAccAlicloudHaVip_update(t *testing.T) {
	var havip vpc.HaVip

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipExists("alicloud_havip.foo", havip),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip.foo", "vswitch_id"),
					resource.TestCheckResourceAttr(
						"alicloud_havip.foo", "description", "tf_testAcc_havip"),
				),
			},
			{
				Config: testAccHaVipUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipExists("alicloud_havip.foo", havip),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip.foo", "vswitch_id"),
					resource.TestCheckResourceAttr(
						"alicloud_havip.foo", "description", "tf_testAcc_havip_update"),
				),
			},
		},
	})
}

func testAccCheckHaVipExists(n string, havip vpc.HaVip) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No HaVip ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		haVipService := HaVipService{client}
		instance, err := haVipService.DescribeHaVip(rs.Primary.ID)
		if err != nil {
			return err
		}
		havip = instance
		return nil
	}
}

func testAccCheckHaVipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	haVipService := HaVipService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_havip" {
			continue
		}
		instance, err := haVipService.DescribeHaVip(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Havip error %#v", err)
		}
		return fmt.Errorf("Havip %s still exist", instance.HaVipId)
	}
	return nil
}

const testAccHaVipConfig = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf_testAcc_havip"
}
 data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}
 resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "tf_testAcc_havip"
}

resource "alicloud_havip" "foo" {
	vswitch_id = "${alicloud_vswitch.foo.id}"
	description = "tf_testAcc_havip"
}
`

const testAccHaVipUpdate = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf_testAcc_havip"
}
 data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}
 resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "tf_testAcc_havip"
}

resource "alicloud_havip" "foo" {
	vswitch_id = "${alicloud_vswitch.foo.id}"
	description = "tf_testAcc_havip_update"
}
`
