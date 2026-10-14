package alicloud

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccAliCloudEcsVsc_basic covers create, in-place update (vsc_name/description),
// tags and import. A pre-existing ECS instance id must be supplied via the
// TF_ECS_VSC_INSTANCE_ID environment variable, because a Vsc is bound to an
// existing instance.
func TestAccAliCloudEcsVsc_basic(t *testing.T) {
	if v := os.Getenv("TF_ECS_VSC_INSTANCE_ID"); v == "" {
		t.Skip("TF_ECS_VSC_INSTANCE_ID is not set, skip TestAccAliCloudEcsVsc_basic")
	}
	instanceId := os.Getenv("TF_ECS_VSC_INSTANCE_ID")
	resourceId := "alicloud_ecs_vsc.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEcsVscDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsVscConfig(instanceId, "tf-testacc-vsc", "desc-before", "Primary", map[string]string{"Created": "TF"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "instance_id", instanceId),
					resource.TestCheckResourceAttr(resourceId, "vsc_name", "tf-testacc-vsc"),
					resource.TestCheckResourceAttr(resourceId, "description", "desc-before"),
					resource.TestCheckResourceAttr(resourceId, "vsc_type", "Primary"),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "TF"),
				),
			},
			{
				Config: testAccEcsVscConfig(instanceId, "tf-testacc-vsc-updated", "desc-after", "Primary", map[string]string{"Created": "TF", "Env": "Acc"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vsc_name", "tf-testacc-vsc-updated"),
					resource.TestCheckResourceAttr(resourceId, "description", "desc-after"),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "2"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func testAccEcsVscConfig(instanceId, name, description, vscType string, tags map[string]string) string {
	return testAccEcsVscTagsConfig(instanceId, name, description, vscType, tags)
}

func testAccEcsVscTagsConfig(instanceId, name, description, vscType string, tags map[string]string) string {
	tagsHCL := ""
	for k, v := range tags {
		tagsHCL += "		" + k + " = \"" + v + "\"\n"
	}
	return `
resource "alicloud_ecs_vsc" "default" {
  instance_id = "` + instanceId + `"
  vsc_name    = "` + name + `"
  description = "` + description + `"
  vsc_type    = "` + vscType + `"
  tags = {
` + tagsHCL + `  }
}
`
}

func testAccCheckEcsVscDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_vsc" {
			continue
		}
		_, err := ecsServiceV2.DescribeEcsVsc(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("Vsc %s still exists", rs.Primary.ID))
	}
	return nil
}
