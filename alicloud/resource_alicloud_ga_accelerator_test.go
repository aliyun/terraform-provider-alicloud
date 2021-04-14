package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaAccelerator_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":            "1",
					"auto_use_coupon": "true",
					"duration":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":            "1",
						"auto_use_coupon": "true",
						"duration":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_use_coupon", "duration"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "accelerator_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "accelerator_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name,
					"description":      "accelerator",
					"spec":             `1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name,
						"description":      "accelerator",
						"spec":             "1",
					}),
				),
			},
		},
	})
}

var AlicloudGaAcceleratorMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence(name string) string {
	return ""
}
