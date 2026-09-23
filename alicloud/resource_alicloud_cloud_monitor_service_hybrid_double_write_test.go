package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudMonitorServiceHybridDoubleWrite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_hybrid_double_write.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudMonitorServiceHybridDoubleWriteMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceHybridDoubleWrite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-chw%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAliCloudCloudMonitorServiceHybridDoubleWriteBasicDependence0)
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
					"source_namespace": "${alicloud_cms_namespace.default.id}",
					"source_user_id":   "${data.alicloud_account.default.id}",
					"namespace":        "${alicloud_cms_namespace.default.id}",
					"user_id":          "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_namespace": CHECKSET,
						"source_user_id":   CHECKSET,
						"namespace":        CHECKSET,
						"user_id":          CHECKSET,
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

var AliCloudCloudMonitorServiceHybridDoubleWriteMap = map[string]string{}

func AliCloudAliCloudCloudMonitorServiceHybridDoubleWriteBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_cms_namespace" "source" {
  		namespace = var.name
	}
	
	resource "alicloud_cms_namespace" "default" {
  		namespace = "${var.name}-source"
	}
`, name)
}
