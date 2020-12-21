package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAlicloudLogProjectTags_basic(t *testing.T) {
	var v *sls.ResourceTag
	resourceId := "alicloud_log_project_tags.default"
	ra := resourceAttrInit(resourceId, logProjectTags)

	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("sls-xuxiaohang-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogProjectTagsDependence)

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
					"project_name":  name,
					"tags": map[string]string{"the-tag":"aliyun-log-go-sdk","the-tag-2":"aliyun log go sdk"},
					"depends_on": []string{"alicloud_log_project.foo"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":        name,
						"tags.%": "2",
						"tags.the-tag" :"aliyun-log-go-sdk",
						"tags.the-tag-2":"aliyun log go sdk",

					}),
				),

			},
		},
	})

}


func resourceLogProjectTagsDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "foo" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

var logProjectTags = map[string]string {
	"project_name": CHECKSET,

}
