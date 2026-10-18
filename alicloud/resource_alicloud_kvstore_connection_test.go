package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kvstore_connection", &resource.Sweeper{
		Name: "alicloud_kvstore_connection",
		F:    testSweepKvstoreConnections,
	})
}

// testSweepKvstoreConnections sweeps KVStore (Redis) public-connection
// sub-resources. A connection is a sub-resource of a KVStore instance, so we
// enumerate instances first, then DescribeDBInstanceNetInfo per instance and
// release the public connection (DBInstanceNetType == "0") if present.
func testSweepKvstoreConnections(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	// 1. List all KVStore instances (Redis + Memcache).
	listAction := "DescribeInstances"
	listReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	instanceIds := make([]string, 0)
	for _, instanceType := range []string{string(KVStoreRedis), string(KVStoreMemcache)} {
		listReq["InstanceType"] = instanceType
		for {
			response, err := client.RpcPost("R-kvstore", "2015-01-01", listAction, nil, listReq, true)
			if err != nil {
				log.Printf("[ERROR] Failed to list KVStore instances for connection sweep: %s", err)
				return nil
			}
			resp, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, listAction, "$.Instances.KVStoreInstance", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				instanceName := fmt.Sprint(item["InstanceName"])
				if protectedSweepVpcName(instanceName) {
					log.Printf("[INFO] Skipping KVStore Instance connection (NODELETING): %s (%s)", instanceName, item["InstanceId"])
					continue
				}
				if !sweepAll() {
					match := false
					for _, prefix := range prefixes {
						if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
							match = true
							break
						}
					}
					if !match {
						continue
					}
				}
				instanceIds = append(instanceIds, fmt.Sprint(item["InstanceId"]))
			}
			if len(result) < PageSizeXLarge {
				break
			}
			listReq["PageNumber"] = listReq["PageNumber"].(int) + 1
		}
	}

	// 2. For each qualifying instance, find its public connection and release it.
	for _, instanceId := range instanceIds {
		netReq := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
		}
		netResp, err := client.RpcPost("R-kvstore", "2015-01-01", "DescribeDBInstanceNetInfo", nil, netReq, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
				continue
			}
			log.Printf("[ERROR] Failed to DescribeDBInstanceNetInfo for instance %s: %s", instanceId, err)
			continue
		}
		netObj, err := jsonpath.Get("$.NetInfoItems.InstanceNetInfo", netResp)
		if err != nil {
			log.Printf("[ERROR] Failed to parse net info for instance %s: %s", instanceId, err)
			continue
		}
		netList, _ := netObj.([]interface{})
		for _, n := range netList {
			info, ok := n.(map[string]interface{})
			if !ok {
				continue
			}
			// DBInstanceNetType == "0" is the public connection (see Read in
			// resource_alicloud_kvstore_connection.go).
			if fmt.Sprint(info["DBInstanceNetType"]) != "0" {
				continue
			}
			connStr := fmt.Sprint(info["ConnectionString"])
			if connStr == "" {
				continue
			}
			log.Printf("[INFO] Deleting KVStore public connection for instance %s: %s", instanceId, connStr)
			delReq := map[string]interface{}{
				"RegionId":              client.RegionId,
				"InstanceId":            instanceId,
				"CurrentConnectionString": connStr,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(1*time.Minute, func() *resource.RetryError {
				_, err = client.RpcPost("R-kvstore", "2015-01-01", "ReleaseInstancePublicConnection", nil, delReq, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound", "InvalidCurrentConnectionString.NotFound"}) {
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
				log.Printf("[ERROR] Failed to release KVStore public connection %s (%s): %s", instanceId, connStr, err)
			}
		}
	}
	return nil
}

func TestAccAliCloudKvstoreConnection_basic(t *testing.T) {
	var v r_kvstore.InstanceNetInfo
	resourceId := "alicloud_kvstore_connection.default"
	ra := resourceAttrInit(resourceId, RedisConnectionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreConnection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreConnectionBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetest",
					"instance_id":              "${alicloud_kvstore_instance.default.id}",
					"port":                     "6370",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              CHECKSET,
						"port":                     "6370",
						"connection_string_prefix": "allocatetest",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_string_prefix"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetestupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "allocatetestupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "6371",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "6371",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetestup",
					"port":                     "6372",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "allocatetestup",
						"port":                     "6372",
					}),
				),
			},
		},
	})
}

var RedisConnectionMap = map[string]string{
	"connection_string": CHECKSET,
}

func KvstoreConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
	  zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	}
	resource "alicloud_kvstore_instance" "default" {
		db_instance_name = "%s"
  		vswitch_id = data.alicloud_vswitches.default.ids.0
		instance_type = "Redis"
		engine_version = "4.0"
		tags = {
			Created = "TF",
			For = "update test",
		}
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_class="redis.master.large.default"
	}
	`, name)
}
