package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCREESyncRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCREESyncRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEESyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
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
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"name":                  name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      defaultRegionToTest,
					"tag_filter":            ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        CHECKSET,
						"name":                  name,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": CHECKSET,
						"target_region_id":      CHECKSET,
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
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCREESyncRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEESyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREESyncRuleBasicDependence1)
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
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"name":                  name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      defaultRegionToTest,
					"tag_filter":            ".*",
					"repo_name":             "${alicloud_cr_ee_repo.source.name}",
					"target_repo_name":      "${alicloud_cr_ee_repo.target.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        CHECKSET,
						"name":                  name,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": CHECKSET,
						"target_region_id":      CHECKSET,
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

	resource "alicloud_cr_ee_namespace" "source" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  		name               = "${var.name}-source"
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_namespace" "target" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.1
  		name               = "${var.name}-target"
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}
`, name)
}

func AliCloudCREESyncRuleBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}

	resource "alicloud_cr_ee_namespace" "source" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  		name               = "${var.name}-source"
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_namespace" "target" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.1
  		name               = "${var.name}-target"
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_repo" "source" {
  		instance_id = alicloud_cr_ee_namespace.source.instance_id
  		namespace   = alicloud_cr_ee_namespace.source.name
  		name        = "${var.name}-source"
  		repo_type   = "PRIVATE"
  		summary     = var.name
	}

	resource "alicloud_cr_ee_repo" "target" {
  		instance_id = alicloud_cr_ee_namespace.target.instance_id
  		namespace   = alicloud_cr_ee_namespace.target.name
  		name        = "${var.name}-target"
  		repo_type   = "PRIVATE"
  		summary     = var.name
	}
`, name)
}
