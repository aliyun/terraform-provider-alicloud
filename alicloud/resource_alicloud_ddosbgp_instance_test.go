package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDdosbgpInstance_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_ddosbgp_instance.default"
	ra := resourceAttrInit(resourceId, ddosbgpInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DdosbgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdosbgpInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdosbgpSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"base_bandwidth":   "20",
					"normal_bandwidth": "100",
					"bandwidth":        "-1",
					"ip_count":         "100",
					"ip_type":          "IPv4",
					"type":             "Enterprise",
					"period":           "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"bandwidth":      "-1",
						"base_bandwidth": "20",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "-update",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":           name,
					"base_bandwidth": "20",
					"ip_count":       "100",
					"ip_type":        "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"base_bandwidth": "20",
						"ip_count":       "100",
						"ip_type":        "IPv4",
					}),
				),
			},
		},
	})
}

func resourceDdosbgpInstanceDependence(name string) string {
	return ``
}

var ddosbgpInstanceBasicMap = map[string]string{
	"ip_count": "100",
	"ip_type":  "IPv4",
}
