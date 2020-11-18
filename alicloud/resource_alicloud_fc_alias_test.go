package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudFCAliasUpdate(t *testing.T) {
	var v *fc.GetAliasOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc-%s-alicloudfcalias-%d-cd", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"alias_name": CHECKSET,
	}
	resourceId := "alicloud_fc_alias.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcAliasConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":    "${alicloud_fc_service.default.name}",
					"alias_name":      "${var.name}",
					"service_version": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":    name,
						"alias_name":      name,
						"service_version": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "this is an alias description.",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "this is an alias description.",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "this is another alias description.",
					"routing_config": []map[string]interface{}{
						{
							"additional_version_weights": map[string]string{
								"1": "0.3",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "this is another alias description.",
						"routing_config.#": "1",
						"routing_config.0.additional_version_weights.1": "0.3",
					}),
				),
			},
		},
	})
}

func resourceFcAliasConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
resource "alicloud_fc_service" "default" {
  name = "${var.name}"
  publish = "true"
}
`, name)
}
