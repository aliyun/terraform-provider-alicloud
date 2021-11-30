package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDTSMigrationInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_migration_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSMigrationInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsMigrationInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSMigrationInstanceBasicDependence0)
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
					"payment_type":                     "PayAsYouGo",
					"source_endpoint_engine_name":      "MySQL",
					"source_endpoint_region":           "cn-hangzhou",
					"destination_endpoint_engine_name": "MySQL",
					"destination_endpoint_region":      "cn-hangzhou",
					"instance_class":                   "small",
					"sync_architecture":                "oneway",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                     "PayAsYouGo",
						"instance_class":                   "small",
						"destination_endpoint_region":      CHECKSET,
						"source_endpoint_engine_name":      "MySQL",
						"source_endpoint_region":           CHECKSET,
						"destination_endpoint_engine_name": "MySQL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
						"From":    "acceptance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "sync_architecture", "compute_unit"},
			},
		},
	})
}

var AlicloudDTSMigrationInstanceMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudDTSMigrationInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
