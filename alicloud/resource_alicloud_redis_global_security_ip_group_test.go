// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
)

func init() {
	resource.AddTestSweepers("alicloud_redis_global_security_ip_group", &resource.Sweeper{
		Name: "alicloud_redis_global_security_ip_group",
		F:    testSweepRedisGlobalSecurityIpGroups,
	})
}

// testSweepRedisGlobalSecurityIpGroups sweeps KVStore (Redis) global security
// IP group sub-resources. Unlike account/backup, a global security IP group is
// a region-level resource not tied to a specific instance, so we list all
// groups directly via DescribeGlobalSecurityIPGroup (GlobalSecurityGroupId
// omitted — the API documents it as optional and returns the full list) and
// delete each via DeleteGlobalSecurityIPGroup. Groups associated with live
// instances must be unbound first; if the cloud rejects such a delete we log
// and continue (non-fatal), leaving the group as a reported leftover.
func testSweepRedisGlobalSecurityIpGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	listReq := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	response, err := client.RpcPost("R-kvstore", "2015-01-01", "DescribeGlobalSecurityIPGroup", nil, listReq, true)
	if err != nil {
		log.Printf("[ERROR] Failed to list KVStore global security IP groups: %s", err)
		return nil
	}
	groupObj, err := jsonpath.Get("$.GlobalSecurityIPGroup", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "DescribeGlobalSecurityIPGroup", "$.GlobalSecurityIPGroup", response)
	}
	groupList, _ := groupObj.([]interface{})
	for _, g := range groupList {
		group, ok := g.(map[string]interface{})
		if !ok {
			continue
		}
		groupName := fmt.Sprint(group["GlobalIgName"])
		groupId := fmt.Sprint(group["GlobalSecurityGroupId"])
		if protectedSweepVpcName(groupName) {
			log.Printf("[INFO] Skipping KVStore global security IP group (NODELETING): %s (%s)", groupName, groupId)
			continue
		}
		if !sweepAll() {
			match := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(groupName), strings.ToLower(prefix)) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		log.Printf("[INFO] Deleting KVStore global security IP group: %s (%s)", groupName, groupId)
		delReq := map[string]interface{}{
			"RegionId":              client.RegionId,
			"GlobalSecurityGroupId": groupId,
			"GlobalIgName":          groupName,
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err = client.RpcPost("R-kvstore", "2015-01-01", "DeleteGlobalSecurityIPGroup", nil, delReq, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"RequiredParam.NotFound"}) {
					return nil
				}
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete KVStore global security IP group %s (%s): %s", groupName, groupId, err)
		}
	}
	return nil
}

// Test Redis GlobalSecurityIpGroup. >>> Resource test cases, automatically generated.
// Case 白名单模板模型测试 12846
func TestAccAliCloudRedisGlobalSecurityIpGroup_basic12846(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisGlobalSecurityIpGroupMap12846)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisGlobalSecurityIpGroupBasicDependence12846)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": "ggn_test_create",
					"global_ip_list":       "192.168.0.1,10.10.10.10,172.16.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": "ggn_test_create",
						"global_ip_list":       "192.168.0.1,10.10.10.10,172.16.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": "ggn_test_update",
					"global_ip_list":       "192.168.0.1,10.10.10.10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": "ggn_test_update",
						"global_ip_list":       "192.168.0.1,10.10.10.10",
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

var AlicloudRedisGlobalSecurityIpGroupMap12846 = map[string]string{}

func AlicloudRedisGlobalSecurityIpGroupBasicDependence12846(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}


`, name)
}

// Test Redis GlobalSecurityIpGroup. <<< Resource test cases, automatically generated.
