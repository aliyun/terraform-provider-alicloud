package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSaeConfigMap_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sae_config_map.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEConfigMapMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeConfigMap")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%ssaeconfigmap%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEConfigMapBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_id": "${alicloud_sae_namespace.default.namespace_id}",
					"name":         "tftestaccname",
					"data":         `{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_id": fmt.Sprintf("%s:configmaptest", os.Getenv("ALICLOUD_REGION")),
						"name":         "tftestaccname",
						"data":         "{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccdescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data": `{\"env.home\":\"/root\",\"env.shell\":\"/bin/sh\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data": "{\"env.home\":\"/root\",\"env.shell\":\"/bin/sh\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAccDesc",
					"data":        `{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccDesc",
						"data":        "{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudSAEConfigMapMap0 = map[string]string{
	"description":  "",
	"namespace_id": fmt.Sprintf("%s:configmaptest", os.Getenv("ALICLOUD_REGION")),
	"name":         "tftestaccname",
}

func AlicloudSAEConfigMapBasicDependence0(name string) string {
	return fmt.Sprintf(` 
resource "alicloud_sae_namespace" "default" {
  namespace_description = "namespace_desc"
  namespace_id = "%s:configmaptest"
  namespace_name = "namespace_name"
}

variable "name" {
  default = "%s"
}
`, os.Getenv("ALICLOUD_REGION"), name)
}
