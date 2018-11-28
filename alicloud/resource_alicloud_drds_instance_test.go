package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_drds_instance", &resource.Sweeper{
		Name: "alicloud_drds_instance",
		F:    testSweepDRDSInstances,
	})
}

func testSweepDRDSInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DrdsSupportedRegions) {
		log.Printf("[INFO] Skipping DRDS Instance unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []drds.Instance
	req := drds.CreateDescribeDrdsInstancesRequest()
	req.RegionId = client.RegionId
	for {
		raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.DescribeDrdsInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving DRDS Instances: %s", err)
		}
		resp, _ := raw.(*drds.DescribeDrdsInstancesResponse)
		if resp == nil || len(resp.Data.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Data.Instance...)

	}

	sweeped := false
	for _, v := range insts {
		name := v.Description
		id := v.DrdsInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping DRDS Instance: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting DRDS Instance: %s (%s)", name, id)
		req := drds.CreateRemoveDrdsInstanceRequest()
		req.DrdsInstanceId = id
		_, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.RemoveDrdsInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DRDS Instance (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudDRDSInstance_Basic(t *testing.T) {
	var instance drds.DescribeDrdsInstanceResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithRegions(t, false, connectivity.DrdsClassicNoSupportedRegions)
		},
		IDRefreshName: "alicloud_drds_instance.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDrdsInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDRDSInstanceExist(
						"alicloud_drds_instance.basic", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.basic",
						"instance_charge_type",
						"PostPaid"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.basic",
						"instance_series",
						"drds.sn1.4c8g"),
					resource.TestCheckResourceAttrSet("alicloud_drds_instance.basic", "zone_id"),

					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.basic",
						"specification",
						"drds.sn1.4c8g.8C16G"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.basic",
						"description",
						"tf-testaccDrdsdatabase_basic"),
				),
			},
		},
	})
}
func TestAccAlicloudDRDSInstance_Vpc(t *testing.T) {
	var instance drds.DescribeDrdsInstanceResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
		},
		IDRefreshName: "alicloud_drds_instance.vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDrdsInstance_Vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDRDSInstanceExist(
						"alicloud_drds_instance.vpc", &instance),
					resource.TestCheckResourceAttrSet(
						"alicloud_drds_instance.vpc",
						"zone_id"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.vpc",
						"instance_charge_type",
						"PostPaid"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.vpc",
						"instance_series",
						"drds.sn1.4c8g"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.vpc",
						"specification",
						"drds.sn1.4c8g.8C16G"),
					resource.TestCheckResourceAttrSet("alicloud_drds_instance.vpc", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_drds_instance.vpc", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.vpc",
						"description",
						"tf-testaccDrdsdatabase_vpc"),
				),
			},
		},
	})
}
func testAccCheckDRDSInstanceExist(n string, instance *drds.DescribeDrdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DRDS Instance ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		drdsService := DrdsService{client}

		req := drds.CreateDescribeDrdsInstanceRequest()
		req.DrdsInstanceId = rs.Primary.ID
		response, err := drdsService.DescribeDrdsInstance(req.DrdsInstanceId)
		if err != nil {
			return err
		} else {
			instance = response
		}
		return nil
	}
}
func testAccCheckDRDSInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_drds_instance" {
			continue
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		drdsService := DrdsService{client}
		req := drds.CreateDescribeDrdsInstanceRequest()
		req.DrdsInstanceId = rs.Primary.ID
		_, err := drdsService.DescribeDrdsInstance(req.DrdsInstanceId)
		if err != nil {
			if NotFoundError(err) {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

const testAccDrdsInstance = `
variable "name" {
	default = "tf-testaccDrdsdatabase_basic"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

variable "instance_series" {
	default = "drds.sn1.4c8g"
}
resource "alicloud_drds_instance" "basic" {
  description = "${var.name}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  instance_series = "${var.instance_series}"
  specification = "drds.sn1.4c8g.8C16G"
}
`
const testAccDrdsInstance_Vpc = `

variable "name" {
	default = "tf-testaccDrdsdatabase_vpc"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

variable "instance_series" {
	default = "drds.sn1.4c8g"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}


resource "alicloud_drds_instance" "vpc" {
  description = "${var.name}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  instance_series = "${var.instance_series}"
  instance_charge_type = "PostPaid"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  vpc_id = "${alicloud_vswitch.foo.vpc_id}"
  specification = "drds.sn1.4c8g.8C16G"
}
`
