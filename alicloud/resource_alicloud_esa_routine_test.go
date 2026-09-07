package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA Routine.
func TestAccAliCloudESARoutine_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine.default"
	ra := resourceAttrInit(resourceId, AliCloudESARoutineMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutine")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARoutineBasicDependence)
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
					"name":             name,
					"description":      "tf-test-routine",
					"code":             "addEventListener('fetch', e => e.respondWith(new Response('v1')))",
					"code_description": "version 1",
					"deploy_env":       "staging",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                CHECKSET,
						"code":                CHECKSET,
						"latest_code_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"code":             "addEventListener('fetch', e => e.respondWith(new Response('v2')))",
					"code_description": "version 2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code":                CHECKSET,
						"latest_code_version": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "code_description", "deploy_env"},
			},
		},
	})
}

var AliCloudESARoutineMap = map[string]string{
	"id":                  CHECKSET,
	"create_time":         CHECKSET,
	"latest_code_version": CHECKSET,
}

func AliCloudESARoutineBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
