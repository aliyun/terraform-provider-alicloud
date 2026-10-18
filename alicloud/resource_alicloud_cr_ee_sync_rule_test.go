package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_ee_sync_rule", &resource.Sweeper{
		Name: "alicloud_cr_ee_sync_rule",
		F:    testSweepCrEESyncRule,
	})
}

// testSweepCrEESyncRule cleans up Container Registry repo sync rules left on
// CR EE instances by acceptance tests. A sync rule can only be deleted from its
// SOURCE instance (sync_direction == "FROM"); rules seen from the target side
// (TO / others) are skipped and left as leftovers because DeleteRepoSyncRule
// refuses them. Rules on instances (or with rule names) marked NODELETING are
// preserved; in non-SWEEP_ALL mode only rules on tf-* named instances are
// removed.
func testSweepCrEESyncRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	prefixes := []string{"tf-testAcc", "tf_testAcc", "tftestAcc"}

	pageNo, pageSize := 1, PageSizeLarge
	for {
		resp, lerr := crService.ListCrEEInstances(pageNo, pageSize)
		if lerr != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s", lerr)
			return nil
		}
		for _, inst := range resp.Instances {
			if protectedSweepVpcName(inst.InstanceName) {
				log.Printf("[INFO] Skipping CR EE sync rule on instance %s (%s): protected (NODELETING)", inst.InstanceId, inst.InstanceName)
				continue
			}
			if !sweepAll() {
				match := false
				for _, p := range prefixes {
					if strings.HasPrefix(strings.ToLower(inst.InstanceName), strings.ToLower(p)) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			listAction := "ListRepoSyncRule"
			listReq := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InstanceId": inst.InstanceId,
				"PageSize":   pageSize,
				"PageNo":     1,
			}
			for {
				wait := incrementalWait(3*time.Second, 3*time.Second)
				var response map[string]interface{}
				rerr := retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcPost("cr", "2018-12-01", listAction, map[string]interface{}{}, listReq, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if rerr != nil {
					log.Printf("[ERROR] Failed to ListRepoSyncRule for instance %s: %s", inst.InstanceId, rerr)
					break
				}
				raw, jerr := jsonpath.Get("$.SyncRules", response)
				if jerr != nil {
					log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", listAction, "$.SyncRules", response)
					break
				}
				result, _ := raw.([]interface{})
				for _, v := range result {
					item := v.(map[string]interface{})
					syncRuleId := fmt.Sprint(item["SyncRuleId"])
					syncRuleName := fmt.Sprint(item["SyncRuleName"])
					syncDirection := fmt.Sprint(item["SyncDirection"])
					if protectedSweepVpcName(syncRuleName) {
						log.Printf("[INFO] Skipping CR EE sync rule %s:%s: protected (NODELETING)", inst.InstanceId, syncRuleName)
						continue
					}
					if syncDirection != "FROM" {
						log.Printf("[INFO] Skipping CR EE sync rule %s:%s: sync_direction=%s (must delete from source instance)", inst.InstanceId, syncRuleName, syncDirection)
						continue
					}
					log.Printf("[INFO] Deleting CR EE sync rule: %s:%s (%s)", inst.InstanceId, syncRuleId, syncRuleName)
					delReq := map[string]interface{}{
						"RegionId":   client.RegionId,
						"InstanceId": inst.InstanceId,
						"SyncRuleId": syncRuleId,
					}
					dwait := incrementalWait(3*time.Second, 3*time.Second)
					if derr := retry.Retry(5*time.Minute, func() *retry.RetryError {
						_, e := client.RpcPost("cr", "2018-12-01", "DeleteRepoSyncRule", map[string]interface{}{}, delReq, true)
						if e != nil {
							if NotFoundError(e) {
								return nil
							}
							if NeedRetry(e) {
								dwait()
								return retry.RetryableError(e)
							}
							return retry.NonRetryableError(e)
						}
						return nil
					}); derr != nil {
						log.Printf("[ERROR] Failed to delete CR EE sync rule %s:%s: %s", inst.InstanceId, syncRuleId, derr)
						continue
					}
				}
				if len(result) < pageSize {
					break
				}
				listReq["PageNo"] = listReq["PageNo"].(int) + 1
			}
		}
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}
	return nil
}

func TestAccAliCloudCrRepoSyncRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_sync_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrRepoSyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"sync_rule_name":        name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      defaultRegionToTest,
					"tag_filter":            ".*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        CHECKSET,
						"sync_rule_name":        name,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": CHECKSET,
						"target_region_id":      CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrRepoSyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_namespace.source.instance_id}",
					"namespace_name":        "${alicloud_cr_ee_namespace.source.name}",
					"sync_rule_name":        name,
					"target_instance_id":    "${alicloud_cr_ee_namespace.target.instance_id}",
					"target_namespace_name": "${alicloud_cr_ee_namespace.target.name}",
					"target_region_id":      defaultRegionToTest,
					"tag_filter":            ".*",
					"repo_name":             "${alicloud_cr_ee_repo.source.name}",
					"target_repo_name":      "${alicloud_cr_ee_repo.target.name}",
					"sync_scope":            "REPO",
					"sync_trigger":          "PASSIVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"namespace_name":        CHECKSET,
						"sync_rule_name":        name,
						"target_instance_id":    CHECKSET,
						"target_namespace_name": CHECKSET,
						"target_region_id":      CHECKSET,
						"tag_filter":            ".*",
						"repo_name":             CHECKSET,
						"target_repo_name":      CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrRepoSyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	ra := resourceAttrInit(resourceId, AliCloudCrRepoSyncRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrRepoSyncRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-sync-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrRepoSyncRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"sync_scope":            "REPO",
					"sync_trigger":          "PASSIVE",
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
	"create_time":       CHECKSET,
	"region_id":         CHECKSET,
	"repo_sync_rule_id": CHECKSET,
	"sync_scope":        CHECKSET,
	"sync_trigger":      CHECKSET,
	"sync_direction":    CHECKSET,
	"rule_id":           CHECKSET,
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
