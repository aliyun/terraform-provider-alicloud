package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudMonitorServiceHybridDoubleWrite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_hybrid_double_write.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudMonitorServiceHybridDoubleWriteMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceHybridDoubleWrite")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-chw%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliCloudCloudMonitorServiceHybridDoubleWriteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"source_namespace": acctest.CHECKSET,
						"source_user_id":   acctest.CHECKSET,
						"namespace":        acctest.CHECKSET,
						"user_id":          acctest.CHECKSET,
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
