package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEssLifecycleHook_basic(t *testing.T) {
	var hook ess.LifecycleHook
	rand1 := acctest.RandIntRange(1000, 999999)
	rand2 := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_lifecycle_hook.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, rand1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssLifecycleHookExists(
						"alicloud_ess_lifecycle_hook.foo", &hook),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"name",
						fmt.Sprintf("tf-testAccEssLifecycleHook-%d", rand1)),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"lifecycle_transition",
						"SCALE_OUT"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"heartbeat_timeout",
						"400"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"notification_metadata",
						"helloworld"),
				),
			},

			{
				Config: testAccEssLifecycleHook_update(EcsInstanceCommonTestCase, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssLifecycleHookExists(
						"alicloud_ess_lifecycle_hook.foo", &hook),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"name",
						fmt.Sprintf("tf-testAccEssLifecycleHook-%d", rand2)),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"lifecycle_transition",
						"SCALE_IN"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"heartbeat_timeout",
						"200"),
					resource.TestCheckResourceAttr(
						"alicloud_ess_lifecycle_hook.foo",
						"notification_metadata",
						"hellojava"),
				),
			},
		},
	})
}

func testAccCheckEssLifecycleHookExists(n string, d *ess.LifecycleHook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS Lifecycle Hook ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		essService := EssService{client}
		attr, err := essService.DescribeLifecycleHookById(rs.Primary.ID)
		log.Printf("[DEBUG] check lifecycle hook %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = attr
		return nil
	}
}

func testAccCheckEssLifecycleHookDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_lifecycle_hook" {
			continue
		}
		if _, err := essService.DescribeLifecycleHookById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("lifecycle hook %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssLifecycleHook(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "bar" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.bar.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "foo"{
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_OUT"
		heartbeat_timeout = 400
		notification_metadata = "helloworld"
	}
	`, common, rand)
}
func testAccEssLifecycleHook_update(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "tf-testAccEssLifecycleHook-%d"
	}
	
	resource "alicloud_vswitch" "bar" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "alicloud_ess_scaling_group" "foo" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}","${alicloud_vswitch.bar.id}"]
	}
	
	resource "alicloud_ess_lifecycle_hook" "foo"{
		scaling_group_id = "${alicloud_ess_scaling_group.foo.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_IN"
		heartbeat_timeout = 200
		notification_metadata = "hellojava"
	}
	`, common, rand)
}
