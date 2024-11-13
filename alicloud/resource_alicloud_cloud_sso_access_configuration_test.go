package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudSSOAccessConfiguration_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessConfigurationBasicDependence0)
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
					"directory_id":              "${data.alicloud_cloud_sso_directories.default.directories.0.id}",
					"access_configuration_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":              CHECKSET,
						"access_configuration_name": name,
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
					"relay_state": "https://cloudsso.console.aliyun.com/test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"relay_state": "https://cloudsso.console.aliyun.com/test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
						{
							"permission_policy_type": "System",
							"permission_policy_name": "AliyunECSFullAccess",
						},
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     "oos",
							"permission_policy_document": `{\"Statement\": [{\"Action\": \"ecs:Get\",\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permission_policies.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     "oos",
							"permission_policy_document": `{\"Statement\": [{\"Action\": \"oss:*\",\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}`,
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCloudSSOAccessConfiguration_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_access_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessConfigurationBasicDependence0)
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
					"directory_id":                     "${data.alicloud_cloud_sso_directories.default.directories.0.id}",
					"access_configuration_name":        name,
					"session_duration":                 "1800",
					"relay_state":                      "https://cloudsso.console.aliyun.com/test1",
					"description":                      name,
					"force_remove_permission_policies": "true",
					"permission_policies": []map[string]interface{}{
						{
							"permission_policy_type": "System",
							"permission_policy_name": "ReadOnlyAccess",
						},
						{
							"permission_policy_type": "System",
							"permission_policy_name": "AliyunECSFullAccess",
						},
						{
							"permission_policy_type":     "Inline",
							"permission_policy_name":     "oos",
							"permission_policy_document": `{\"Statement\": [{\"Action\": \"ecs:Get\",\"Effect\": \"Allow\",\"Resource\": \"*\"}],\"Version\": \"1\"}`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":              CHECKSET,
						"access_configuration_name": name,
						"session_duration":          "1800",
						"relay_state":               "https://cloudsso.console.aliyun.com/test1",
						"description":               name,
						"permission_policies.#":     "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_remove_permission_policies"},
			},
		},
	})
}

var AliCloudCloudSSOAccessConfigurationMap0 = map[string]string{
	"session_duration":        CHECKSET,
	"access_configuration_id": CHECKSET,
}

func AliCloudCloudSSOAccessConfigurationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cloud_sso_directories" "default" {
	}
`, name)
}
