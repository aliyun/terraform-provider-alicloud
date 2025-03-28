package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA EdgeContainerApp. >>> Resource test cases, automatically generated.
// Case resource_EdgeContainerApp_test
func TestAccAliCloudESAEdgeContainerAppresource_EdgeContainerApp_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_edge_container_app.default"
	ra := resourceAttrInit(resourceId, AliCloudESAEdgeContainerAppresource_EdgeContainerApp_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaEdgeContainerApp")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAEdgeContainerApp%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAEdgeContainerAppresource_EdgeContainerApp_testBasicDependence)
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
					"health_check_fail_times": "5",
					"service_port":            "80",
					"target_port":             "3000",
					"health_check_interval":   "5",
					"health_check_host":       "example.com",
					"health_check_uri":        "/",
					"health_check_timeout":    "3",
					"health_check_succ_times": "2",
					"remarks":                 "无备注信息1",
					"health_check_method":     "HEAD",
					"edge_container_app_name": "terraform-app-002",
					"health_check_port":       "80",
					"health_check_http_code":  "http_2xx",
					"health_check_type":       "l7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAEdgeContainerAppresource_EdgeContainerApp_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAEdgeContainerAppresource_EdgeContainerApp_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ESA EdgeContainerApp. <<< Resource test cases, automatically generated.
