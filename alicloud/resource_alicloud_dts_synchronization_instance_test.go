package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDTSSynchronizationInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_synchronization_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSSynchronizationInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsSynchronizationInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSSynchronizationInstanceBasicDependence0)
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                     "PayAsYouGo",
						"source_endpoint_engine_name":      "MySQL",
						"source_endpoint_region":           "cn-hangzhou",
						"destination_endpoint_engine_name": "MySQL",
						"destination_endpoint_region":      "cn-hangzhou",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"database_count", "status", "quantity", "sync_architecture", "auto_start", "compute_unit", "period", "used_time", "auto_pay", "order_type", "synchronization_direction"},
			},
		},
	})
}

var AlicloudDTSSynchronizationInstanceMap0 = map[string]string{
	"sync_architecture":         NOSET,
	"auto_start":                NOSET,
	"compute_unit":              NOSET,
	"period":                    NOSET,
	"used_time":                 NOSET,
	"auto_pay":                  NOSET,
	"order_type":                NOSET,
	"synchronization_direction": NOSET,
	"database_count":            NOSET,
	"status":                    NOSET,
	"quantity":                  NOSET,
}

func AlicloudDTSSynchronizationInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
