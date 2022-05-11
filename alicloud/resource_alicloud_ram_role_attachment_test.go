package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudRAMRoleAttachment_basic(t *testing.T) {
	var instanceA ecs.Instance
	var instanceB ecs.Instance
	var v *ecs.DescribeInstanceRamRoleResponse
	resourceId := "alicloud_ram_role_attachment.default"
	ra := resourceAttrInit(resourceId, ramRoleAttachmentMap)
	serviceFuncRam := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	serviceFuncEcs := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFuncRam)
	rcInstanceA := resourceCheckInit("alicloud_instance.default.0", &instanceA, serviceFuncEcs)
	rcInstanceB := resourceCheckInit("alicloud_instance.default.1", &instanceB, serviceFuncEcs)

	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamRoleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRoleAttachmentConfig(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcInstanceA.checkResourceExists(),
					rcInstanceB.checkResourceExists(),
					testAccCheck(nil),
				),
			},
		},
	})
}

var ramRoleAttachmentMap = map[string]string{
	"role_name":      CHECKSET,
	"instance_ids.#": "2",
}

func testAccRamRoleAttachmentConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAcc%sRamRoleAttachmentConfig-%d"
	}

	resource "alicloud_instance" "default" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"

		# series III
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		count = 2

		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		security_groups = ["${alicloud_security_group.default.id}"]
	}

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}

	resource "alicloud_ram_role_attachment" "default" {
	  role_name = "${alicloud_ram_role.default.name}"
	  instance_ids = "${alicloud_instance.default.*.id}"
	}`, common, defaultRegionToTest, rand)
}

func testAccCheckRamRoleAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_role_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ecs.CreateDescribeInstanceRamRoleRequest()
		request.InstanceIds = strings.Split(rs.Primary.ID, ":")[1]

		for {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeInstanceRamRole(request)
			})
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				continue
			}
			if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
				break
			}
			if err == nil {
				response, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
				if len(response.InstanceRamRoleSets.InstanceRamRoleSet) > 0 {
					for _, v := range response.InstanceRamRoleSets.InstanceRamRoleSet {
						if v.RamRoleName != "" {
							return WrapError(fmt.Errorf("Attach %s still exists.", rs.Primary.ID))
						}
					}
				}
				break
			}
			return WrapError(err)
		}
	}
	return nil
}
