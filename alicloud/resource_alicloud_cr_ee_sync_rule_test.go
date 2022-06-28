package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCREESyncRule_Basic(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	var v *cr_ee.SyncRulesItem
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEESyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEESyncRuleConfigDependence)

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
					"instance_id":           "${alicloud_cr_ee_namespace.source_ns.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source_ns.name}",
					"name":                  "${var.name}",
					"target_region_id":      region,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target_ns.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target_ns.name}",
					"tag_filter":            ".*",
					"depends_on": []string{
						"alicloud_cr_ee_repo.source_repo",
						"alicloud_cr_ee_repo.target_repo",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        name,
						"name":                  name,
						"target_region_id":      region,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": name,
						"tag_filter":            ".*",
						"rule_id":               CHECKSET,
						"sync_direction":        "FROM",
						"sync_scope":            "NAMESPACE",
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

func resourceCrEESyncRuleConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_cr_ee_instances" "default" {}

resource "alicloud_cr_ee_namespace" "source_ns" {
	instance_id = data.alicloud_cr_ee_instances.default.ids.0
	name = "${var.name}"
	auto_create	= true
	default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_namespace" "target_ns" {
	instance_id = data.alicloud_cr_ee_instances.default.ids.1
	name = "${var.name}"
	auto_create	= true
	default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "source_repo" {
	instance_id = alicloud_cr_ee_namespace.source_ns.instance_id
	namespace = alicloud_cr_ee_namespace.source_ns.name
	name = var.name
	summary = "test"
	repo_type = "PRIVATE"
	detail = "test"
}

resource "alicloud_cr_ee_repo" "target_repo" {
	instance_id = "${alicloud_cr_ee_namespace.target_ns.instance_id}"
	namespace = "${alicloud_cr_ee_namespace.target_ns.name}"
	name = "${var.name}"
	summary = "test"
	repo_type = "PRIVATE"
	detail = "test"
}
`, name)
}
