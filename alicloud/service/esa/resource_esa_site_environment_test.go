package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA EsaSiteEnvironment. >>> Resource test cases, automatically generated.
// Case Test EsaSiteEnvironment 8508
func TestAccAliCloudEsaSiteEnvironment_basic8508(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_environment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaSiteEnvironmentMap8508)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaSiteEnvironment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("Environment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaSiteEnvironmentBasicDependence8508)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":               "${data.alicloud_esa_sites.default.sites.0.id}",
					"environment_name":      name,
					"rule":                  "(ip.src in {1.1.1.1})",
					"next_environment_name": "Production",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":               acctest.CHECKSET,
						"environment_name":      name,
						"rule":                  "(ip.src in {1.1.1.1})",
						"next_environment_name": "Production",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(ip.src in {1.1.1.2})",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": "(ip.src in {1.1.1.2})",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"next_environment_name"},
			},
		},
	})
}

func TestAccAliCloudEsaSiteEnvironment_basic8508_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_environment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaSiteEnvironmentMap8508)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaSiteEnvironment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("Environment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaSiteEnvironmentBasicDependence8508)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":               "${data.alicloud_esa_sites.default.sites.0.id}",
					"environment_name":      name,
					"rule":                  "(ip.src in {1.1.1.1})",
					"next_environment_name": "Production",
					"read_only":             true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":               acctest.CHECKSET,
						"environment_name":      name,
						"rule":                  "(ip.src in {1.1.1.1})",
						"next_environment_name": "Production",
						"read_only":             "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"next_environment_name"},
			},
		},
	})
}

var AliCloudEsaSiteEnvironmentMap8508 = map[string]string{
	"create_time": acctest.CHECKSET,
	"update_time": acctest.CHECKSET,
}

func AliCloudEsaSiteEnvironmentBasicDependence8508(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}
`, name)
}

// Test ESA EsaSiteEnvironment. <<< Resource test cases, automatically generated.
