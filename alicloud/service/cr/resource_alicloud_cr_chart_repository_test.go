package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCRChartRepository_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_chart_repository.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCRChartRepositoryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrChartRepository")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCRChartRepositoryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_namespace_name": "${alicloud_cr_chart_namespace.default.namespace_name}",
					"instance_id":         "${alicloud_cr_chart_namespace.default.instance_id}",
					"repo_name":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_namespace_name": acctest.CHECKSET,
						"instance_id":         acctest.CHECKSET,
						"repo_name":           name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PUBLIC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": "summary_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": "summary_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PRIVATE",
					"summary":   "Summary_update_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PRIVATE",
						"summary":   "Summary_update_all",
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

var AlicloudCRChartRepositoryMap0 = map[string]string{
	"repo_name":           acctest.CHECKSET,
	"repo_namespace_name": acctest.CHECKSET,
}

func AlicloudCRChartRepositoryBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cr_ee_instances" "default" {}

resource "alicloud_cr_chart_namespace" "default" {
	instance_id        = data.alicloud_cr_ee_instances.default.ids.0
	namespace_name       = var.name
}


`, name)
}

// lintignore: R001
