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
					"sync_rule_name":        name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      acctest.DefaultRegionToTest,
					"tag_filter":            ".*",
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
