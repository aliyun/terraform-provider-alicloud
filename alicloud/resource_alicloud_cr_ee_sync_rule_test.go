package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCREESyncRule_basic0(t *testing.T) {
	var v *cr_ee.SyncRulesItem
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCREESyncRuleMap0)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEESyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc-creesr-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREESyncRuleBasicDependence0)
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
					"instance_id":           "${alicloud_cr_ee_repo.source_repo.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_repo.source_repo.namespace}",
					"target_instance_id":    "${alicloud_cr_ee_repo.target_repo.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_repo.target_repo.namespace}",
					"target_region_id":      defaultRegionToTest,
					"name":                  name,
					"tag_filter":            ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        CHECKSET,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": CHECKSET,
						"target_region_id":      CHECKSET,
						"name":                  name,
						"tag_filter":            ".*",
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

func TestAccAliCloudCREESyncRule_basic0_twin(t *testing.T) {
	var v *cr_ee.SyncRulesItem
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCREESyncRuleMap0)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEESyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc-creesr-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREESyncRuleBasicDependence0)
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
					"instance_id":           "${alicloud_cr_ee_repo.source_repo.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_repo.source_repo.namespace}",
					"target_instance_id":    "${alicloud_cr_ee_repo.target_repo.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_repo.target_repo.namespace}",
					"target_region_id":      defaultRegionToTest,
					"name":                  name,
					"tag_filter":            ".*",
					"repo_name":             "${alicloud_cr_ee_repo.source_repo.name}",
					"target_repo_name":      "${alicloud_cr_ee_repo.target_repo.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        CHECKSET,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": CHECKSET,
						"target_region_id":      CHECKSET,
						"name":                  name,
						"tag_filter":            ".*",
						"repo_name":             CHECKSET,
						"target_repo_name":      CHECKSET,
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

var AliCloudCREESyncRuleMap0 = map[string]string{
	"rule_id":        CHECKSET,
	"sync_direction": CHECKSET,
	"sync_scope":     CHECKSET,
}

func AliCloudCREESyncRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}

	resource "alicloud_cr_ee_namespace" "source_ns" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  		name               = var.name
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_namespace" "target_ns" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.1
  		name               = var.name
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_repo" "source_repo" {
  		instance_id = alicloud_cr_ee_namespace.source_ns.instance_id
  		namespace   = alicloud_cr_ee_namespace.source_ns.name
  		name        = var.name
  		summary     = "test"
  		repo_type   = "PRIVATE"
  		detail      = var.name
	}

	resource "alicloud_cr_ee_repo" "target_repo" {
  		instance_id = alicloud_cr_ee_namespace.target_ns.instance_id
  		namespace   = alicloud_cr_ee_namespace.target_ns.name
  		name        = var.name
  		summary     = "test"
  		repo_type   = "PRIVATE"
  		detail      = var.name
	}
`, name)
}
