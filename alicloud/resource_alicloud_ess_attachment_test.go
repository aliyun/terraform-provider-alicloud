package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEssAttachment_basic(t *testing.T) {
	var sg ess.ScalingGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_attachment.attach",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAttachmentConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 99999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alicloud_ess_attachment.attach", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_ess_attachment.attach",
						"instance_ids.#", "2"),
				),
			},
		},
	})
}

func testAccCheckEssAttachmentExists(n string, d *ess.ScalingGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS attachment ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		group, err := essService.DescribeScalingGroup(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}

		instances, err := essService.DescribeScalingInstances(rs.Primary.ID, "", make([]string, 0), string(Attached))

		if err != nil {
			return WrapError(err)
		}

		if len(instances) < 1 {
			return WrapError(Error("Scaling instances not found"))
		}

		*d = group
		return nil
	}
}

func testAccCheckEssAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}

		_, err := essService.DescribeScalingGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidScalingGroupIdNotFound) {
				continue
			}
			return WrapError(err)
		}

		instances, err := essService.DescribeScalingInstances(rs.Primary.ID, "", make([]string, 0), string(Attached))

		if err != nil && !IsExceptedError(err, InvalidScalingGroupIdNotFound) {
			return WrapError(err)
		}

		if len(instances) > 0 {
			return WrapError(fmt.Errorf("There are still ECS instances in the scaling group."))
		}
	}

	return nil
}

func testAccEssAttachmentConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAttachmentConfig-%d"
	}
	data "alicloud_instance_types" "special" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count    = 2
		memory_size       = 4
		instance_type_family = "ecs.sn1ne"
	}

	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 0
		max_size = 2
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

	resource "alicloud_ess_scaling_configuration" "foo" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"

		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
		security_group_id = "${alicloud_security_group.default.id}"
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_instance" "instance" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.special.instance_types.0.id}"
		count = 2
		security_groups = ["${alicloud_security_group.default.id}"]
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = "10"
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_ess_attachment" "attach" {
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		instance_ids = ["${alicloud_instance.instance.0.id}", "${alicloud_instance.instance.1.id}"]
		force = true
	}
	`, common, rand)
}
