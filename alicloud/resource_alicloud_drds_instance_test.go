package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

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
		if err == nil && response != nil && response.Data.Status != "5" {
			return fmt.Errorf("error! DRDS instance still exists")
		}
	}
	return nil
}

const testAccDrdsInstance = `
provider "alicloud" {
	region = "cn-hangzhou"
}
resource "alicloud_drds_instance" "basic" {
  provider = "alicloud"
  description = "drds basic"
  type = "PRIVATE"
  zone_id = "cn-hangzhou-e"
  specification = "drds.sn1.4c8g.8C16G"
  pay_type = "drdsPost"
  instance_series = "drds.sn1.4c8g"
}
`
const testAccDrdsInstance_Vpc = `
provider "alicloud" {
	region = "cn-hangzhou"
}
resource "alicloud_drds_instance" "vpc" {
  provider = "alicloud"
  description = "drds vpc"
  type = "PRIVATE"
  zone_id = "cn-hangzhou-e"
  specification = "drds.sn1.4c8g.16C32G"
  pay_type = "drdsPost"
  vswitch_id = "vsw-bp1rfn58rx73af8oswzye"
  instance_series = "drds.sn1.4c8g"
}
`
