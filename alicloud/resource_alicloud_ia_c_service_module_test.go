package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test IaCService Module. >>> Resource test cases, automatically generated.
// Case Module lifecycle test
func TestAccAliCloudIaCServiceModule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ia_c_service_module.default"
	ra := resourceAttrInit(resourceId, AlicloudIaCServiceModuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &IaCServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeIaCServiceModule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%siacmodule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudIaCServiceModuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"module_name":      name,
					"source":           "Registry",
					"source_path":      "alibaba/security-group:2.4.1",
					"version_strategy": "Manual",
					"description":      "tf-testacc module",
					"tags": []map[string]interface{}{
						{
							"tag_key":   "Created",
							"tag_value": "TF",
						},
						{
							"tag_key":   "For",
							"tag_value": "Test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"module_name":      name,
						"source":           "Registry",
						"source_path":      "alibaba/security-group:2.4.1",
						"version_strategy": "Manual",
						"description":      "tf-testacc module",
						"tags.#":           "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"module_name": name + "_update",
					"description": "tf-testacc module update",
					"source_path": "terraform-alicloud-modules/mongodb:3.0.0",
					"state_path":  "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform.tfstate",
					"tags": []map[string]interface{}{
						{
							"tag_key":   "Created-update",
							"tag_value": "TF-update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"module_name": name + "_update",
						"description": "tf-testacc module update",
						"source_path": "terraform-alicloud-modules/mongodb:3.0.0",
						"state_path":  "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform.tfstate",
						"tags.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_strategy": "SourcePathUpdated",
					"state_path":       "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform-update.tfstate",
					"tags":             REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_strategy": "SourcePathUpdated",
						"state_path":       "oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform-update.tfstate",
						"tags.#":           "0",
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

var AlicloudIaCServiceModuleMap0 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"output_path": CHECKSET,
}

func AlicloudIaCServiceModuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test IaCService Module. <<< Resource test cases, automatically generated.
