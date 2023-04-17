package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudCmsHybridDoubleWrite_basic2811(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_hybrid_double_write.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsHybridDoubleWriteMap2811)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsHybridDoubleWrite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccmshybriddoublewrite%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsHybridDoubleWriteBasicDependence2811)
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
					"user_id":          "${data.alicloud_account.default.id}",
					"source_namespace": "${alicloud_cms_namespace.default.0.namespace}",
					"namespace":        "${alicloud_cms_namespace.default.1.namespace}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":          CHECKSET,
						"source_namespace": CHECKSET,
						"namespace":        CHECKSET,
					}),
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

var AlicloudCmsHybridDoubleWriteMap2811 = map[string]string{}

func AlicloudCmsHybridDoubleWriteBasicDependence2811(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {}

resource "alicloud_cms_namespace" "default" {
	count = 2
	description = var.name
	namespace = "${var.name}-${count.index}"
	specification = "cms.s1.large"
}

`, name)
}
