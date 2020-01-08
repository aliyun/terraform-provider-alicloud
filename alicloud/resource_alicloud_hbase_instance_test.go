package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const resourceHBaseConfigClassic = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
}
`

const resourceHBaseConfigClassicName = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic_change_name"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
}
`

const resourceHBaseConfigVpc = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.vswitches.0.id}"
  cold_storage_size = 0
}
`

const resourceHBaseConfigVpcName = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.vswitches.0.id}"
  cold_storage_size = 0
}
`

const resourceHBaseConfigMultiInstance = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  count = 2
  name = "tf_testAccHBase_multi"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
}
`
const resourceHBaseConfigMultiInstanceChangeName = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  count = 2
  name = "tf_testAccHBase_multi_change_name"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
}
`

func TestAccAlicloudHBaseInstanceClassic(t *testing.T) {
	var instance hbase.DescribeInstanceResponse

	resourceId := "alicloud_hbase_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf_testAccHBase_classic",
						"engine_version":     "2.0",
						"core_instance_type": "hbase.n1.large",
						"core_disk_type":     "cloud_efficiency",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceHBaseConfigClassicName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_classic_change_name",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudHBaseInstanceVpc(t *testing.T) {
	var instance hbase.DescribeInstanceResponse

	resourceId := "alicloud_hbase_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf_testAccHBase_vpc",
						"engine_version":     "2.0",
						"core_instance_type": "hbase.n1.large",
						"core_disk_type":     "cloud_efficiency",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceHBaseConfigVpcName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_vpc_change_name",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudHBaseInstanceMultiInstance(t *testing.T) {
	var instance hbase.DescribeInstanceResponse

	resourceId := "alicloud_hbase_instance.default.1"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigMultiInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf_testAccHBase_multi",
						"engine_version":     "2.0",
						"core_instance_type": "hbase.n1.large",
						"core_disk_type":     "cloud_efficiency",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				Config: resourceHBaseConfigMultiInstanceChangeName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_multi_change_name",
					}),
				),
			},
		},
	})
}
