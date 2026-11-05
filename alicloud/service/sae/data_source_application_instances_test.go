package sae_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSaeApplicationInstancesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// ids filter, matching
			{
				Config: testAccApplicationInstancesConfig(name, fmt.Sprintf(`ids = ["${data.alicloud_sae_application_instances.image.instances.0.id}"]`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_restarts"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_health_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.package_version"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.image_url"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.debug_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.main_container_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.timestamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.create_time_stamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.finish_time_stamp"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.eip", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.unhealthy_message", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.traffic_status", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.sidecar_containers_status.#", "0"),
				),
			},
			// ids filter, no match
			{
				Config: testAccApplicationInstancesConfig(name, fmt.Sprintf(`ids = ["${data.alicloud_sae_application_instances.image.instances.0.id}_fake"]`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "0"),
				),
			},
			// instance_id filter, matching
			{
				Config: testAccApplicationInstancesConfig(name, fmt.Sprintf(`instance_id = "${data.alicloud_sae_application_instances.image.instances.0.id}"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_restarts"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_health_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.package_version"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.image_url"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.debug_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.main_container_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.timestamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.create_time_stamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.finish_time_stamp"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.eip", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.unhealthy_message", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.traffic_status", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.sidecar_containers_status.#", "0"),
				),
			},
			// instance_id filter, no match
			{
				Config: testAccApplicationInstancesConfig(name, fmt.Sprintf(`instance_id = "${data.alicloud_sae_application_instances.image.instances.0.id}_fake"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "0"),
				),
			},
			// both filters, matching
			{
				Config: testAccApplicationInstancesConfig(name, fmt.Sprintf(`ids = ["${data.alicloud_sae_application_instances.image.instances.0.id}"]
  instance_id = "${data.alicloud_sae_application_instances.image.instances.0.id}"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.group_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_restarts"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_container_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.instance_health_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.package_version"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.image_url"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.debug_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.main_container_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.timestamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.create_time_stamp"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.default", "instances.0.finish_time_stamp"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.eip", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.unhealthy_message", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.traffic_status", ""),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.sidecar_containers_status.#", "0"),
				),
			},
			// both filters, no match
			{
				Config: testAccApplicationInstancesConfig(name, fmt.Sprintf(`ids = ["${data.alicloud_sae_application_instances.image.instances.0.id}_fake"]
  instance_id = "${data.alicloud_sae_application_instances.image.instances.0.id}_fake"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "0"),
				),
			},
		},
	})
}

func testAccApplicationInstancesConfig(name string, filter string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_sae_applications" "image" {
  app_name = "default-NODELETING-image"
}

data "alicloud_sae_application_groups" "image" {
  app_id = data.alicloud_sae_applications.image.applications.0.id
}

data "alicloud_sae_application_instances" "image" {
  app_id   = data.alicloud_sae_application_groups.image.app_id
  group_id = data.alicloud_sae_application_groups.image.groups.0.id
}

data "alicloud_sae_application_instances" "default" {
  app_id   = data.alicloud_sae_application_groups.image.app_id
  group_id = data.alicloud_sae_application_groups.image.groups.0.id
  %s
}
`, name, filter)
}
