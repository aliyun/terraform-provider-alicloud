package alicloud

import (
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbAttachment_basic(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlb")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
					}),
				),
			},
			{
				Config: testAccSlbAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "70",
					}),
				),
			},
			{
				Config: testAccSlbAttachmentUpdateInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccSlbAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbAttachment_multi(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_attachment.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlb")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbAttachmentMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbAttachment_classic_basic(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlb")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbAttachmentClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
					}),
				),
			},
			{
				Config: testAccSlbAttachmentClassicUpdateWeight,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "70",
					}),
				),
			},
			{
				Config: testAccSlbAttachmentClassicUpdateInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccSlbAttachmentClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
					}),
				),
			},
		},
	})
}

const testAccSlbAttachment = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}
data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_attachment" "default" {
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.id}"]
	weight = 90
}

`
const testAccSlbAttachmentUpdate = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_attachment" "default" {
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.id}"]
	weight = 70
}
`

const testAccSlbAttachmentUpdateInstance = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}
data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	count = 2
	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_attachment" "default" {
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.0.id}","${alicloud_instance.default.1.id}"]
	weight = 70
}
`

const testAccSlbAttachmentMulti = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_attachment" "default" {
	count = 10
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.id}"]
	weight = 90
}
`

const testAccSlbAttachmentClassic = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
}

resource "alicloud_slb_attachment" "default" {
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.id}"]
	weight = 90
}
`

const testAccSlbAttachmentClassicUpdateWeight = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
}

resource "alicloud_slb_attachment" "default" {
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.id}"]
	weight = 70
}
`

const testAccSlbAttachmentClassicUpdateInstance = `
variable "name" {
	default = "tf-testAccSlbAttachment"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"

}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"
	count = 2
	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
	name = "${var.name}"
}

resource "alicloud_slb_attachment" "default" {
	load_balancer_id = "${alicloud_slb.default.id}"
	instance_ids = ["${alicloud_instance.default.0.id}","${alicloud_instance.default.1.id}"]
	weight = 70
}
`
