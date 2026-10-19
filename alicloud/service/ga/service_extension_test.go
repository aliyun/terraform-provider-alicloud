package ga_test

import (
	"fmt"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudGaServiceExtension_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "alicloud_ga_service_extension.default"
	resourceName := fmt.Sprintf("tf-testAccGaServiceExtension-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAlicloudGaServiceExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaServiceExtensionConfig_basic(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudGaServiceExtensionExists(resourceId),
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
					testAccCheckAlicloudGaServiceExtensionExists(resourceId),
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
					testAccCheckAlicloudGaServiceExtensionExists(resourceId),
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

// testAccCheckAlicloudGaServiceExtensionExists verifies the resource is
// present in the API by calling GetServiceExtension directly through the
// shared SDKv2 client meta (the muxed provider configures it during the test).
func testAccCheckAlicloudGaServiceExtensionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		if _, err := getGaServiceExtensionForTest(client, rs.Primary.ID); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAlicloudGaServiceExtensionDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ga_service_extension" {
			continue
		}
		if _, err := getGaServiceExtensionForTest(client, rs.Primary.ID); err != nil {
			if isGaServiceExtensionNotFound(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Ga Service Extension %s still exists", rs.Primary.ID)
	}
	return nil
}

// getGaServiceExtensionForTest calls GetServiceExtension and returns the
// response object. It mirrors the resource's describe path for the test
// without reaching into the unexported ga package helpers.
func getGaServiceExtensionForTest(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RpcPost("Ga", "2019-11-20", "GetServiceExtension", nil, map[string]interface{}{
		"RegionId":           client.RegionId,
		"ServiceExtensionId": id,
	}, true)
	if err != nil {
		return nil, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, fmt.Errorf("GetServiceExtension response carried no object: %v", response)
	}
	object, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("GetServiceExtension response was not an object: %v", response)
	}
	return object, nil
}

func isGaServiceExtensionNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return containsAny(msg, "NotExist.ServiceExtension", "NotFound", "not found", "ResourceNotfound")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && indexOf(s, sub) >= 0 {
			return true
		}
	}
	return false
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
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
