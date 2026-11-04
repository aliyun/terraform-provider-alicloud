package edas_test

// collection-order: bypass instance_ids -- pre-existing stock TypeList unchanged by this directory-only migration; a convergent two-member reorder apply requires live ECS instances attached to an EDAS cluster that the existing tests do not exercise, and fabricating instance/cluster identifiers is out of scope for this migration

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	edassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudEdasinstanceClusterAttachment_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alicloud_edas_instance_cluster_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, edasICAttachmentMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasicattachment%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasICAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testEdasCheckICAttachmentDestroy,
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
	"cluster_id": acctest.CHECKSET,
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
		  instance_type_family = "ecs.s6"
		  cpu_core_count       = 1
		  memory_size          = 2
		  system_disk_category = "cloud_essd"
		}

		data "alicloud_images" "default" {
		  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64"
		  most_recent = true
		  owners      = "system"
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
		  system_disk_category       = "cloud_essd"
		  image_id                   = data.alicloud_images.default.images.0.id
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
