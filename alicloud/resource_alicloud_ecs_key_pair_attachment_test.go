package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckKeyPairAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_key_pair_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		instanceIds := rs.Primary.Attributes["instance_ids"]

		for _, inst := range instanceIds {
			response, err := ecsService.DescribeInstance(string(inst))
			if err != nil {
				return err
			}

			if response.KeyPairName != "" {
				return fmt.Errorf("Error Key Pair Attachment still exist")
			}

		}
	}

	return nil
}

func TestAccAlicloudECSKeyPairAttachmentBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPairAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsKeyPairAttachmentBasicDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var AlicloudEcsKeyPairAttachmentMap = map[string]string{
	"key_pair_name":  CHECKSET,
	"instance_ids.#": "1",
}

func AlicloudEcsKeyPairAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "${var.name}"
}
resource "alicloud_ecs_key_pair_attachment" "default" {
  key_pair_name = "${alicloud_ecs_key_pair.default.id}"
  instance_ids = [alicloud_instance.default.id]
}
`, name)
}
