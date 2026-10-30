package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA RoutineCodeDeployment.
func TestAccAliCloudESARoutineCodeDeployment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_code_deployment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESARoutineCodeDeploymentMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRoutineCodeDeployment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESARoutineCodeDeploymentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"deployment_id":              acctest.CHECKSET,
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
	"deployment_id": acctest.CHECKSET,
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
