package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA RoutineCodeDeployment.
func TestAccAliCloudESARoutineCodeDeployment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_code_deployment.default"
	ra := resourceAttrInit(resourceId, AliCloudESARoutineCodeDeploymentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineCodeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARoutineCodeDeploymentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"routine_name": "${alicloud_esa_routine.default.name}",
					"env":          "staging",
					"strategy":     "percentage",
					"code_versions": []map[string]interface{}{
						{
							"code_version": "${alicloud_esa_routine.default.latest_code_version}",
							"percentage":   "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env":                        "staging",
						"code_versions.#":            "1",
						"code_versions.0.percentage": "100",
						"deployment_id":              CHECKSET,
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

var AliCloudESARoutineCodeDeploymentMap = map[string]string{
	"deployment_id": CHECKSET,
}

func AliCloudESARoutineCodeDeploymentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_esa_routine" "default" {
    name             = var.name
    description      = "tf-test-routine"
    code             = "addEventListener('fetch', e => e.respondWith(new Response('hello')))"
    code_description = "version 1"
}
`, name)
}
