package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AliCloudESARoutineCodeDeploymentMap = map[string]string{
	"deployment_id": CHECKSET,
}

func AliCloudESARoutineCodeDeploymentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// TestAccAliCloudESARoutineCodeDeployment_basic covers CreateRoutineCodeDeployment for a
// single-version percentage deployment, and updates the deployed code version by rolling
// the routine to a new committed code version.
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
	name := fmt.Sprintf("tftestacc%sesadeploy%d", defaultRegionToTest, rand)

	codeV1 := esaRoutineWriteCodeFixture(t, name+"-v1", "addEventListener('fetch', e => e.respondWith(new Response('deploy-v1')))")
	codeV2 := esaRoutineWriteCodeFixture(t, name+"-v2", "addEventListener('fetch', e => e.respondWith(new Response('deploy-v2')))")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccESARoutineCodeDeploymentConfig(name, codeV1, "code version v1", "staging", "percentage"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"routine_name":                 name,
						"env":                          "staging",
						"strategy":                     "percentage",
						"code_versions.#":              "1",
						"code_versions.0.code_version": CHECKSET,
						"code_versions.0.percentage":   "100",
					}),
				),
			},
			{
				Config: testAccESARoutineCodeDeploymentConfig(name, codeV2, "code version v2", "staging", "percentage"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code_versions.#":              "1",
						"code_versions.0.code_version": CHECKSET,
						"code_versions.0.percentage":   "100",
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

func testAccESARoutineCodeDeploymentConfig(name, filename, codeDescription, env, strategy string) string {
	return fmt.Sprintf(`
resource "alicloud_esa_routine" "default" {
  name             = "%[1]s"
  description      = "tf-testacc esa routine deployment"
  filename         = "%[2]s"
  code_description = "%[3]s"
}

resource "alicloud_esa_routine_code_deployment" "default" {
  routine_name = alicloud_esa_routine.default.name
  env          = "%[4]s"
  strategy     = "%[5]s"
  code_versions {
    code_version = alicloud_esa_routine.default.latest_code_version
    percentage   = 100
  }
}
`, name, filename, codeDescription, env, strategy)
}
