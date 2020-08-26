package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEdasinstanceClusterAttachment_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alicloud_edas_instance_cluster_attachment.default"
	ra := resourceAttrInit(resourceId, edasICAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasicattachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasICAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckICAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${alicloud_edas_cluster.default.id}",
					"instance_ids": []string{"${alicloud_instance.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasICAttachmentMap = map[string]string{
	"cluster_id": CHECKSET,
}

func testEdasCheckICAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasICAttachmentDependence(name string) string {
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
		  vswitch_id                 = alicloud_vswitch.default.id
		  security_groups            = ["${alicloud_security_group.default.id}"]
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${alicloud_vpc.default.id}"
		}
		`, name)
}
