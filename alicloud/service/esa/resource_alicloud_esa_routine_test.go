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

// Test ESA Routine.
func TestAccAliCloudESARoutine_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESARoutineMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRoutine")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESARoutineBasicDependence)
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
					"name":             name,
					"description":      "tf-test-routine",
					"code":             "addEventListener('fetch', e => e.respondWith(new Response('v1')))",
					"code_description": "version 1",
					"deploy_env":       "staging",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                acctest.CHECKSET,
						"code":                acctest.CHECKSET,
						"latest_code_version": acctest.CHECKSET,
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
						"code":                acctest.CHECKSET,
						"latest_code_version": acctest.CHECKSET,
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
	"id":                  acctest.CHECKSET,
	"create_time":         acctest.CHECKSET,
	"latest_code_version": acctest.CHECKSET,
}

func AliCloudESARoutineBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
