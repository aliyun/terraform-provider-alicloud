package sae_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSAEApplicationInstancesDataSource(t *testing.T) {
	rand := sdkacctest.RandIntRange(1, 1000)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	outputFile := filepath.Join(os.TempDir(), fmt.Sprintf("tfacc-sae-application-instances-%d.json", rand))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSAEApplicationInstancesDataSourceConfig(rand, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.default", "instances.0.image_url", fmt.Sprintf("registry-vpc.%s.aliyuncs.com/sae-demo-image/consumer:1.0", region)),
					resource.TestCheckResourceAttrPair("data.alicloud_sae_application_instances.default", "instances.0.vswitch_id", "data.alicloud_vswitches.default", "ids.0"),
					resource.TestCheckResourceAttrPair("data.alicloud_sae_application_instances.default", "ids.0", "data.alicloud_sae_application_instances.default", "instances.0.instance_id"),
				),
			},
			{
				Config: testAccSAEApplicationInstancesDataSourceConfigByGroup(rand, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.by_group", "instances.#", "1"),
					resource.TestCheckResourceAttrPair("data.alicloud_sae_application_instances.by_group", "instances.0.instance_id", "data.alicloud_sae_application_instances.default", "instances.0.instance_id"),
					resource.TestCheckResourceAttrPair("data.alicloud_sae_application_instances.by_group", "instances.0.group_id", "data.alicloud_sae_application_instances.default", "instances.0.group_id"),
				),
			},
			{
				Config: testAccSAEApplicationInstancesDataSourceConfigNoMatch(rand, region, outputFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.no_match", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_sae_application_instances.no_match", "ids.#", "0"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_application_instances.no_match", "output_file"),
				),
			},
		},
	})
}

func testAccSAEApplicationInstancesDataSourceBase(rand int, region string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-sae-example-%d"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = "%s:tftfexample%d"
  namespace_name        = var.name
}
resource "alicloud_sae_application" "default" {
  app_description    = var.name
  app_name           = var.name
  namespace_id       = alicloud_sae_namespace.default.namespace_id
  image_url          = "registry-vpc.%s.aliyuncs.com/sae-demo-image/consumer:1.0"
  package_type       = "Image"
  vswitch_id         = data.alicloud_vswitches.default.ids.0
  vpc_id             = data.alicloud_vpcs.default.ids.0
  replicas           = "1"
  cpu                = "500"
  memory             = "2048"
  micro_registration = "0"
}
`, rand, region, rand, region)
}

func testAccSAEApplicationInstancesDataSourceConfig(rand int, region string) string {
	return fmt.Sprintf(`%s
data "alicloud_sae_application_instances" "default" {
  application_id = alicloud_sae_application.default.id
}
`, testAccSAEApplicationInstancesDataSourceBase(rand, region))
}

func testAccSAEApplicationInstancesDataSourceConfigByGroup(rand int, region string) string {
	return fmt.Sprintf(`%s
data "alicloud_sae_application_instances" "default" {
  application_id = alicloud_sae_application.default.id
}
data "alicloud_sae_application_instances" "by_group" {
  application_id = alicloud_sae_application.default.id
  group_id       = data.alicloud_sae_application_instances.default.instances.0.group_id
}
`, testAccSAEApplicationInstancesDataSourceBase(rand, region))
}

func testAccSAEApplicationInstancesDataSourceConfigNoMatch(rand int, region, outputFile string) string {
	return fmt.Sprintf(`%s
data "alicloud_sae_application_instances" "no_match" {
  application_id = alicloud_sae_application.default.id
  ids            = ["fake-instance-id"]
  output_file    = "%s"
}
`, testAccSAEApplicationInstancesDataSourceBase(rand, region), outputFile)
}
