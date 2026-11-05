package ecs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The account under test may hold no reserved instance at all, and purchasing
// one is out of scope for an acceptance test. The non-empty result mapping is
// covered by the data source unit tests; here the empty result, the complete
// filter request and the argument validation are exercised against the live API.
func TestAccAliCloudReservedInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	resourceId := "data.alicloud_reserved_instances.default"
	name := fmt.Sprintf("tf-testAcc-ri-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf("  reserved_instance_name = %q\n", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock("  ids = [\"ri-does-not-exist\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf(`  reserved_instance_name = %q
  instance_type          = "ecs.g6.large"
  instance_type_family   = "ecs.g6"
  scope                  = "Region"
  offering_type          = "No Upfront"
  allocation_type        = "Normal"
  lock_reason            = "financial"
  status                 = ["Active"]
  page_number            = 1
  page_size              = 10
  tags = {
    env = "tf-not-matching"
  }
`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: reservedInstancesDataBlock(fmt.Sprintf(`  reserved_instance_name = %q
  scope                  = "Zone"
  zone_id                = data.alicloud_zones.default.zones.0.id
`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				// zone_id without scope keeps the documented compatibility
				// behaviour: the read requests Scope = Zone instead of failing.
				Config: reservedInstancesDataBlock("  zone_id = data.alicloud_zones.default.zones.0.id\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      reservedInstancesDataBlock("  status = [\"NotAStatus\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      reservedInstancesDataBlock("  scope = \"InvalidScope\"\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      reservedInstancesDataBlock("  scope = \"Zone\"\n"),
				ExpectError: regexp.MustCompile("must be set when"),
			},
			{
				Config:      reservedInstancesDataBlock("  tags = { \"env\" = \"\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must not be empty`),
			},
			{
				// A null value must be rejected even though the key itself is
				// valid; it must not be serialized as an empty value.
				Config:      reservedInstancesDataBlock("  tags = { \"env\" = null }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must be set`),
			},
			{
				Config:      reservedInstancesDataBlock("  tags = { \"acs:env\" = \"test\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Key|must not start with`),
			},
			{
				// The destroy phase re-reads the last configuration, so a valid
				// step must follow the expected validation failures.
				Config: reservedInstancesDataBlock(fmt.Sprintf("  reserved_instance_name = %q\n", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instances.#", "0"),
				),
			},
		},
	})
}

func reservedInstancesDataBlock(arguments string) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {}

data "alicloud_reserved_instances" "default" {
%s}
`, arguments)
}
