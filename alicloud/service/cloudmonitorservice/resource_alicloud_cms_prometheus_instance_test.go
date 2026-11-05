// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cms PrometheusInstance. >>> Resource test cases, automatically generated.
// Case promInstance 8018
func TestAccAliCloudCmsPrometheusInstance_basic8018(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_prometheus_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsPrometheusInstanceMap8018)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsPrometheusInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsPrometheusInstanceBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_instance_name": name,
					"workspace":                "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_instance_name": name,
						"workspace":                acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"archive_duration": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"archive_duration": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_free_read_policy": "1.1.1.1",
					"enable_auth_free_read": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_free_read_policy": "1.1.1.1",
						"enable_auth_free_read": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_free_read_policy": "2.2.2.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_free_read_policy": "2.2.2.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_free_write_policy": "2.2.2.2",
					"enable_auth_free_write": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_free_write_policy": "2.2.2.2",
						"enable_auth_free_write": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_free_write_policy": "1.1.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_free_write_policy": "1.1.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auth_free_read": "false",
					"auth_free_read_policy": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auth_free_read": "false",
						"auth_free_read_policy": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auth_free_write": "false",
					"auth_free_write_policy": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auth_free_write": "false",
						"auth_free_write_policy": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_duration": "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_duration": "180",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudCmsPrometheusInstance_basic8018_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_prometheus_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCmsPrometheusInstanceMap8018)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsPrometheusInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCmsPrometheusInstanceBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_instance_name": name,
					"workspace":                "${alicloud_cms_workspace.default.id}",
					"archive_duration":         "60",
					"auth_free_read_policy":    "1.1.1.1",
					"enable_auth_free_read":    "true",
					"auth_free_write_policy":   "2.2.2.2",
					"enable_auth_free_write":   "true",
					"storage_duration":         "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_instance_name": name,
						"workspace":                acctest.CHECKSET,
						"archive_duration":         "60",
						"auth_free_read_policy":    "1.1.1.1",
						"enable_auth_free_read":    "true",
						"auth_free_write_policy":   "2.2.2.2",
						"enable_auth_free_write":   "true",
						"storage_duration":         "180",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudCmsPrometheusInstanceMap8018 = map[string]string{
	"payment_type": acctest.CHECKSET,
	"create_time":  acctest.CHECKSET,
	"region_id":    acctest.CHECKSET,
}

func AliCloudCmsPrometheusInstanceBasicDependence8018(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
`, name)
}

// Test Cms PrometheusInstance. <<< Resource test cases, automatically generated.
