package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCrRepoSyncRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrRepoSyncRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// step0: NAMESPACE scope, target = the second instance (ids.1).
				// Cover the new write-side fields: link_id, namespace_name_filter,
				// priority, repo_name_filter and target_user_id (current account).
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"sync_rule_name":        name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      acctest.DefaultRegionToTest,
					"tag_filter":            ".*",
					"sync_trigger":          "PASSIVE",
					"target_user_id":        "${data.alicloud_caller_identity.current.account_id}",
					"link_id":               "stl-test-1",
					"priority":              2,
					"namespace_name_filter": "ns.*",
					"repo_name_filter":      "repo.*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           acctest.CHECKSET,
						"namespace_name":        acctest.CHECKSET,
						"sync_rule_name":        name,
						"target_instance_id":    acctest.CHECKSET,
						"target_namespace_name": acctest.CHECKSET,
						"target_region_id":      acctest.CHECKSET,
						"tag_filter":            ".*",
						"sync_trigger":          "PASSIVE",
						"priority":              "2",
						"link_id":               "stl-test-1",
						"namespace_name_filter": "ns.*",
						"repo_name_filter":      "repo.*",
					}),
				),
			},
			{
				// step1: switch to REPO scope and retarget to a namespace inside the
				// source instance (ids.0) to exercise modification of every newly
				// updatable field: sync_scope, repo_name, target_repo_name,
				// target_instance_id, target_namespace_name, target_region_id,
				// link_id, priority, namespace_name_filter, repo_name_filter,
				// sync_rule_name, tag_filter and sync_trigger.
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"sync_rule_name":        name + "-update",
					"target_instance_id":    "${alicloud_cr_ee_namespace.target_in_source.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target_in_source.name}",
					"target_region_id":      "cn-shanghai",
					"tag_filter":            "v.*",
					"sync_trigger":          "INITIATIVE",
					"target_user_id":        "${data.alicloud_caller_identity.current.account_id}",
					"link_id":               "stl-test-2",
					"priority":              4,
					"namespace_name_filter": "ns2.*",
					"repo_name_filter":      "repo2.*",
					"sync_scope":            "REPO",
					"repo_name":             "${alicloud_cr_ee_repo.source.name}",
					"target_repo_name":      "${alicloud_cr_ee_repo.target_in_source.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           acctest.CHECKSET,
						"namespace_name":        acctest.CHECKSET,
						"sync_rule_name":        name + "-update",
						"target_instance_id":    acctest.CHECKSET,
						"target_namespace_name": acctest.CHECKSET,
						"target_region_id":      acctest.CHECKSET,
						"tag_filter":            "v.*",
						"sync_trigger":          "INITIATIVE",
						"sync_scope":            "REPO",
						"repo_name":             acctest.CHECKSET,
						"target_repo_name":      acctest.CHECKSET,
						"priority":              "4",
						"link_id":               "stl-test-2",
						"namespace_name_filter": "ns2.*",
						"repo_name_filter":      "repo2.*",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_user_id"},
			},
		},
	})
}

func TestAccAliCloudCrRepoSyncRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrRepoSyncRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"sync_rule_name":        name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      acctest.DefaultRegionToTest,
					"tag_filter":            ".*",
					"repo_name":             "${alicloud_cr_ee_repo.source.name}",
					"target_repo_name":      "${alicloud_cr_ee_repo.target.name}",
					"sync_scope":            "REPO",
					"sync_trigger":          "PASSIVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           acctest.CHECKSET,
						"namespace_name":        acctest.CHECKSET,
						"sync_rule_name":        name,
						"target_instance_id":    acctest.CHECKSET,
						"target_namespace_name": acctest.CHECKSET,
						"target_region_id":      acctest.CHECKSET,
						"tag_filter":            ".*",
						"repo_name":             acctest.CHECKSET,
						"target_repo_name":      acctest.CHECKSET,
						"sync_scope":            "REPO",
						"sync_trigger":          "PASSIVE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_user_id"},
			},
		},
	})
}

func TestAccAliCloudCrRepoSyncRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrRepoSyncRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"name":                  name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      acctest.DefaultRegionToTest,
					"tag_filter":            ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           acctest.CHECKSET,
						"namespace_name":        acctest.CHECKSET,
						"name":                  name,
						"target_instance_id":    acctest.CHECKSET,
						"target_namespace_name": acctest.CHECKSET,
						"target_region_id":      acctest.CHECKSET,
						"tag_filter":            ".*",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_user_id"},
			},
		},
	})
}

func TestAccAliCloudCrRepoSyncRule_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrRepoSyncRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"name":                  name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      acctest.DefaultRegionToTest,
					"tag_filter":            ".*",
					"repo_name":             "${alicloud_cr_ee_repo.source.name}",
					"target_repo_name":      "${alicloud_cr_ee_repo.target.name}",
					"sync_scope":            "REPO",
					"sync_trigger":          "PASSIVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           acctest.CHECKSET,
						"namespace_name":        acctest.CHECKSET,
						"name":                  name,
						"target_instance_id":    acctest.CHECKSET,
						"target_namespace_name": acctest.CHECKSET,
						"target_region_id":      acctest.CHECKSET,
						"tag_filter":            ".*",
						"repo_name":             acctest.CHECKSET,
						"target_repo_name":      acctest.CHECKSET,
						"sync_scope":            "REPO",
						"sync_trigger":          "PASSIVE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_user_id"},
			},
		},
	})
}

var AliCloudCrRepoSyncRuleMap0 = map[string]string{
	"create_time":       acctest.CHECKSET,
	"region_id":         acctest.CHECKSET,
	"repo_sync_rule_id": acctest.CHECKSET,
	"sync_scope":        acctest.CHECKSET,
	"sync_trigger":      acctest.CHECKSET,
	"sync_direction":    acctest.CHECKSET,
	"rule_id":           acctest.CHECKSET,
	"cross_user":        acctest.CHECKSET,
	"priority":          acctest.CHECKSET,
}

func AliCloudCrRepoSyncRuleBasicDependence0(name string) string {
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

func AliCloudCrRepoSyncRuleBasicDependence1(name string) string {
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

func AliCloudCrRepoSyncRuleBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cr_ee_instances" "default" {
	}

	data "alicloud_caller_identity" "current" {
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

	resource "alicloud_cr_ee_namespace" "target_in_source" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  		name               = "${var.name}-target-in-source"
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

	resource "alicloud_cr_ee_repo" "target_in_source" {
  		instance_id = alicloud_cr_ee_namespace.target_in_source.instance_id
  		namespace   = alicloud_cr_ee_namespace.target_in_source.name
  		name        = "${var.name}-target-in-source"
  		repo_type   = "PRIVATE"
  		summary     = var.name
	}
`, name)
}
