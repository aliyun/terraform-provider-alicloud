package fc_test

import (
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	fcsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/fc"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccAlicloudFCAliasUpdate(t *testing.T) {
	var v *fc.GetAliasOutput
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc-%s-alicloudfcalias-%d-cd", acctest.DefaultRegionToTest, rand)
	var basicMap = map[string]string{
		"alias_name": acctest.CHECKSET,
	}
	resourceId := "alicloud_fc_alias.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return fcsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceFcAliasConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
