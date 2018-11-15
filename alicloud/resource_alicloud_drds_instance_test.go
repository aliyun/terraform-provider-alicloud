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
			testAccPreCheck(t)
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
						"alicloud_drds_instance.foo",
						"type",
						"1"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"pay_type",
						"Postpaid"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"instance_series",
						"drds.sn1.4c8g"),

					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"zone_id",
						"cn-hangzhou-e"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"specification",
						"drds.sn1.4c8g.8C16G"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"description",
						"drds basic"),
				),
			},
		},
	})
}
func TestAccAlicloudDRDSInstance_Vpc(t *testing.T) {
	var instance drds.DescribeDrdsInstanceResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
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
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"zone_id",
						"cn-hangzhou-e"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"type",
						"1"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"pay_type",
						"Postpaid"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"instance_series",
						"drds.sn1.4c8g"),
					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"specification",
						"drds.sn1.4c8g.8C16G"),

					resource.TestCheckResourceAttr(
						"alicloud_drds_instance.foo",
						"vswitch_id",
						"vsw-wz94tq5g4qaj4ri2rhonn"),
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
		if err == nil && response != nil && response.Data.DrdsInstanceId != "" {
			instance = response
			return nil
		}
		return fmt.Errorf("error finding DRDS instance %#v", rs.Primary.ID)
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
		response, err := drdsService.DescribeDrdsInstance(req.DrdsInstanceId)
		if err == nil && response != nil {
			return fmt.Errorf("error! DRDS instance still exists : %s", err)
		}
	}
	return nil
}

const testAccDrdsInstance = `
variable "zone_id" {
	default = "cn-hangzhou-e"
}

variable "instance_series" {
	default = "drds.sn1.4c8g"
}
resource "alicloud_drds_instance" "basic" {
  description = "drds basic"
  type = "1"
  zone_id = "${var.zone_id}"
  pay_type = "Postpaid"
  instance_series = "${var.instance_series}"
  specification = "drds.sn1.4c8g.8C16G"
  vswitch_id="vsw-wz94tq5g4qaj4ri2rhonn"
}
`
const testAccDrdsInstance_Vpc = `

variable "zone_id" {
	default = "cn-hangzhou-f"
}

variable "instance_series" {
	default = "drds.sn1.4c8g"
}

variable "vswitch_id"{
	default = "vsw-bp1jlu3swk8rq2yoi40ey"
}


resource "alicloud_drds_instance" "vpc" {
  provider = "alicloud"
  description = "drds vpc"
  type = "PRIVATE"
  zone_id = "${var.zone_id}"
  instance_series = "${var.instance_series}"
  pay_type = "Postpaid"
  vswitch_id = "${var.vswitch_id}"
}
`
