package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudDtsInstance_basic1170(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudDtsInstanceMap1170)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDtsInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDtsInstanceBasicDependence1170)
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
					"type":                             "migration",
					"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"payment_type":                     "PayAsYouGo",
					"instance_class":                   "large",
					"source_endpoint_engine_name":      "MySQL",
					"source_region":                    os.Getenv("ALICLOUD_REGION"),
					"destination_endpoint_engine_name": "MySQL",
					"compute_unit":                     "2",
					"database_count":                   "1",
					"destination_region":               os.Getenv("ALICLOUD_REGION"),
					"sync_architecture":                "oneway",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                             "migration",
						"resource_group_id":                CHECKSET,
						"payment_type":                     "PayAsYouGo",
						"instance_class":                   "large",
						"source_endpoint_engine_name":      "MySQL",
						"source_region":                    CHECKSET,
						"destination_endpoint_engine_name": "MySQL",
						"compute_unit":                     "2",
						"database_count":                   "1",
						"destination_region":               CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"compute_unit", "database_count", "used_time", "sync_architecture"},
			},
		},
	})
}

func TestAccAliCloudDtsInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudDtsInstanceMap1170)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDtsInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDtsInstanceBasicDependence1170)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":                             "sync",
					"resource_group_id":                "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"payment_type":                     "Subscription",
					"instance_class":                   "large",
					"source_endpoint_engine_name":      "MySQL",
					"source_region":                    os.Getenv("ALICLOUD_REGION"),
					"destination_endpoint_engine_name": "MySQL",
					"compute_unit":                     "2",
					"database_count":                   "1",
					"destination_region":               os.Getenv("ALICLOUD_REGION"),
					"auto_start":                       "true",
					"du":                               "30",
					"period":                           "Month",
					"sync_architecture":                "bidirectional",
					"used_time":                        "1",
					"fee_type":                         "ONLY_CONFIGURATION_FEE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                             "sync",
						"resource_group_id":                CHECKSET,
						"payment_type":                     "Subscription",
						"instance_class":                   "large",
						"source_endpoint_engine_name":      "MySQL",
						"source_region":                    CHECKSET,
						"destination_endpoint_engine_name": "MySQL",
						"compute_unit":                     "2",
						"database_count":                   "1",
						"destination_region":               CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"compute_unit", "database_count", "used_time", "auto_start", "sync_architecture", "period", "du", "fee_type"},
			},
		},
	})
}

var AliCloudDtsInstanceMap1170 = map[string]string{}

func AliCloudDtsInstanceBasicDependence1170(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}
