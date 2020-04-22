package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEdasSlbAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_slb_bind.default"

	ra := resourceAttrInit(resourceId, edasSLBAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasslbattachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasSLBAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckSLBAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id": "${alicloud_edas_application.default.id}",
					"slb_id": "${alicloud_slb.default.id}",
					"slb_ip": "${alicloud_slb.default.address}",
					"type":   "${alicloud_slb.default.address_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

var edasSLBAttachmentMap = map[string]string{
	"app_id": CHECKSET,
	"slb_id": CHECKSET,
	"slb_ip": CHECKSET,
	"type":   CHECKSET,
}

func testEdasCheckSLBAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasSLBAttachmentDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_instance_types" "default" {
		  cpu_core_count    = 1
		  memory_size       = 2
		}

		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "alicloud_vswitch" "default" {
		  vpc_id            = "${alicloud_vpc.default.id}"
		  cidr_block        = "172.16.0.0/24"
		  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
		  name              = "${var.name}"
		}

		resource "alicloud_security_group" "default" {
		  name   = "${var.name}"
		  vpc_id = "${alicloud_vpc.default.id}"
		}

		resource "alicloud_security_group_rule" "default" {
			type = "ingress"
			ip_protocol = "tcp"
			nic_type = "intranet"
			policy = "accept"
			port_range = "22/22"
			priority = 1
			security_group_id = "${alicloud_security_group.default.id}"
			cidr_ip = "172.16.0.0/24"
		}

		resource "alicloud_instance" "default" {
		  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
		  system_disk_category       = "cloud_efficiency"
		  image_id                   = "centos_7_06_64_20G_alibase_20190711.vhd"
		  instance_name              = "${var.name}"
		  vswitch_id                 = "${alicloud_vswitch.default.id}"
		  security_groups            = ["${alicloud_security_group.default.id}"]
		  internet_max_bandwidth_out = 10
		}
		
		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${alicloud_vpc.default.id}"
		}
		
		resource "alicloud_edas_instance_cluster_attachment" "default" {
		  cluster_id = "${alicloud_edas_cluster.default.id}"
		  instance_ids = ["${alicloud_instance.default.id}"]
		}
		
		resource "alicloud_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${alicloud_edas_cluster.default.id}"
		  package_type = "JAR"
		  ecu_info = ["${alicloud_edas_instance_cluster_attachment.default.ecu_map[alicloud_instance.default.id]}"]
		}

		resource "alicloud_slb" "default" {
		  name = "${var.name}"
		}
		`, name)
}
