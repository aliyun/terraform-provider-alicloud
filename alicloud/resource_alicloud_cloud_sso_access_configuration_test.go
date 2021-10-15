package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudSSOAccessConfiguration_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudSSOAccessConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSOAccessConfigurationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_configuration_name": name,
					"directory_id":              "${local.directory_id}",
					"session_duration":          "1200",
					"description":               name,
					"relay_state":               "https://cloudsso.console.aliyun.com/test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_configuration_name": name,
						"directory_id":              CHECKSET,
						"session_duration":          "1200",
						"description":               name,
						"relay_state":               "https://cloudsso.console.aliyun.com/test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"session_duration": "1800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"session_duration": "1800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"relay_state": "https://cloudsso.console.aliyun.com/test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"relay_state": "https://cloudsso.console.aliyun.com/test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     name,
							"permission_policy_document": "\\n{\\n  \\\"Statement\\\": [\\n    {\\n      \\\"Action\\\": \\\"oss:*\\\",\\n      \\\"Effect\\\": \\\"Allow\\\",\\n      \\\"Resource\\\": \\\"*\\\"\\n    }\\n  ],\\n  \\\"Version\\\": \\\"1\\\"\\n}\\n                        ",
						},
						{
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     name,
							"permission_policy_document": "\\n{\\n  \\\"Statement\\\": [\\n    {\\n      \\\"Action\\\": \\\"oss:*\\\",\\n      \\\"Effect\\\": \\\"Allow\\\",\\n      \\\"Resource\\\": \\\"*\\\"\\n    }\\n  ],\\n  \\\"Version\\\": \\\"1\\\"\\n}\\n                        ",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"session_duration": "1200",
					"description":      name,
					"relay_state":      "https://cloudsso.console.aliyun.com/test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"session_duration": "1200",
						"description":      name,
						"relay_state":      "https://cloudsso.console.aliyun.com/test1",
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

var AlicloudCloudSSOAccessConfigurationMap0 = map[string]string{
	"access_configuration_id": CHECKSET,
	"permission_policies.#":   CHECKSET,
	"session_duration":        CHECKSET,
	"directory_id":            CHECKSET,
}

func AlicloudCloudSSOAccessConfigurationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cloud_sso_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name    = var.name
}
locals{
  directory_id =  length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
`, name)
}
