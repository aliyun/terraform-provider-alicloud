package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNASLifecyclePolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_lifecycle_policy.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASLifecyclePolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasLifecyclePolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snaslifecyclepolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASLifecyclePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_policy_name": "${var.name}",
					"file_system_id":        "${alicloud_nas_file_system.default.id}",
					"lifecycle_rule_name":   "DEFAULT_ATIME_14",
					"storage_type":          "InfrequentAccess",
					"paths":                 []string{"/"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_policy_name": name,
						"file_system_id":        acctest.CHECKSET,
						"lifecycle_rule_name":   "DEFAULT_ATIME_14",
						"storage_type":          "InfrequentAccess",
						"paths.#":               "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule_name": "DEFAULT_ATIME_30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule_name": "DEFAULT_ATIME_30",
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

var AlicloudNASLifecyclePolicyMap0 = map[string]string{
	"file_system_id":        acctest.CHECKSET,
	"storage_type":          acctest.CHECKSET,
	"lifecycle_policy_name": acctest.CHECKSET,
	"paths.#":               acctest.CHECKSET,
}

func AlicloudNASLifecyclePolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}
`, name)
}

// lintignore: R001
