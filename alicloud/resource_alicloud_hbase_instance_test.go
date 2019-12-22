package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const resourceHBaseConfigClassic = `
resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic"
  zone_id = "cn-shenzhen-b"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "Postpaid"
  is_cold_storage = "false"
  security_ip_list = ["127.0.0.1", "127.0.0.2"]
}
`

const resourceHBaseConfigClassicName = `
resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic_change_name"
  zone_id = "cn-shenzhen-b"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "Postpaid"
  is_cold_storage = "false"
  security_ip_list = ["127.0.0.1", "127.0.0.2"]
}
`

const resourceHBaseConfigVpc = `
resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "Postpaid"
  vswitch_id = "vsw-wz9iqvmkdua0svi31ox61"
  is_cold_storage = "false"
}
`

const resourceHBaseConfigVpcName = `
resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "Postpaid"
  vswitch_id = "vsw-wz9iqvmkdua0svi31ox61"
  is_cold_storage = "false"
  security_ip_list = ["127.0.0.1", "127.0.0.2"]
}
`

const resourceHBaseConfigMultiInstance = `
resource "alicloud_hbase_instance" "default" {
  count = 2
  name = "tf_testAccHBase_multi"
  zone_id = "cn-shenzhen-b"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "Postpaid"
  is_cold_storage = "false"
  security_ip_list = ["127.0.0.1", "127.0.0.2"]
}
`
const resourceHBaseConfigMultiInstanceChangeName = `
resource "alicloud_hbase_instance" "default" {
  count = 2
  name = "tf_testAccHBase_multi_change_name"
  zone_id = "cn-shenzhen-b"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "Postpaid"
  is_cold_storage = "false"
  security_ip_list = ["127.0.0.1", "127.0.0.2"]
}
`

func TestAccAlicloudHBaseInstanceClassic(t *testing.T) {
	var instance hbase.Instance

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

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "hbase",
						"engine_version": "2.0",
						"net_type":       "CLASSIC",
						"status":         "ACTIVATION",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
			},
			// hbase 参数一个一个修改
			{
				Config: resourceHBaseConfigClassicName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_classic_change_name",
					}),
				),
			},
			// hbase 参数一起修改
		},
	})
}

func TestAccAlicloudHBaseInstanceVpc(t *testing.T) {
	var instance hbase.Instance

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

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "hbase",
						"engine_version": "2.0",
						"net_type":       "VPC",
						"status":         "ACTIVATION",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
			},
			// change value
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
	var instance hbase.Instance

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

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigMultiInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "hbase",
						"engine_version": "2.0",
						"net_type":       "CLASSIC",
						"status":         "ACTIVATION",
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
