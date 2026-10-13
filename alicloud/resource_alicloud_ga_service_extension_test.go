package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudGaServiceExtension_basic(t *testing.T) {
	var extension map[string]interface{}
	rand := acctest.RandInt()
	resourceId := "alicloud_ga_service_extension.default"
	resourceName := fmt.Sprintf("tf-testAccGaServiceExtension-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudGaServiceExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaServiceExtensionConfig_basic(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudGaServiceExtensionExists(resourceId, &extension),
					resource.TestCheckResourceAttr(resourceId, "name", resourceName),
					resource.TestCheckResourceAttr(resourceId, "description", "test description for service extension"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "tfTest"),
					resource.TestCheckResourceAttr(resourceId, "tags.For", "acceptance"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"components", "resources", "create_time", "update_time", "type", "state", "resource_group_id", "tags"},
			},
			{
				Config: testAccGaServiceExtensionConfig_update(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudGaServiceExtensionExists(resourceId, &extension),
					resource.TestCheckResourceAttr(resourceId, "name", resourceName+"-update"),
					resource.TestCheckResourceAttr(resourceId, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceId, "tags.Updated", "true"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"components", "resources", "create_time", "update_time", "type", "state", "resource_group_id", "tags"},
			},
			{
				Config: testAccGaServiceExtensionConfig_notags(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudGaServiceExtensionExists(resourceId, &extension),
					resource.TestCheckResourceAttr(resourceId, "name", resourceName+"-notags"),
					resource.TestCheckResourceAttr(resourceId, "description", "description after clearing tags"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"components", "resources", "create_time", "update_time", "type", "state", "resource_group_id", "tags"},
			},
		},
	})
}

func testAccCheckAlicloudGaServiceExtensionExists(n string, extension *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		gaService := GaService{client}
		resp, err := gaService.DescribeGaServiceExtension(rs.Primary.ID)
		if err != nil {
			return err
		}
		*extension = resp
		return nil
	}
}

func testAccCheckAlicloudGaServiceExtensionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	gaService := GaService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ga_service_extension" {
			continue
		}
		_, err := gaService.DescribeGaServiceExtension(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Ga Service Extension %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccGaServiceExtensionConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ga_service_extension" "default" {
  name        = "%[1]s"
  description = "test description for service extension"
  tags = {
    Created = "tfTest"
    For     = "acceptance"
  }
}`, name)
}

func testAccGaServiceExtensionConfig_update(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ga_service_extension" "default" {
  name        = "%[1]s-update"
  description = "updated description"
  tags = {
    Updated = "true"
  }
}`, name)
}

func testAccGaServiceExtensionConfig_notags(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ga_service_extension" "default" {
  name        = "%[1]s-notags"
  description = "description after clearing tags"
}`, name)
}
