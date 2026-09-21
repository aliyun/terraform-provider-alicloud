package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_direct_mail_dedicated_ip_pool",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_dedicated_ip_pool",
			F:    testSweepDirectMailDedicatedIpPool,
		})
}

func testSweepDirectMailDedicatedIpPool(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.DmSupportRegions)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DedicatedIpPoolList"
	request := map[string]interface{}{
		"PageIndex": 1,
		"PageSize":  PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch DirectMail Dedicated Ip Pool: %s", err)
			return nil
		}
		v, err := jsonpath.Get("$.IpPools", response)
		if err != nil {
			return nil
		}
		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Name"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DirectMail Dedicated Ip Pool: %s", name)
				continue
			}
			action := "DedicatedIpPoolDelete"
			deleteRequest := map[string]interface{}{
				"Id": item["Id"],
			}
			_, err = client.RpcPost("Dm", "2015-11-23", action, nil, deleteRequest, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DirectMail Dedicated Ip Pool (%s): %s", name, err)
			} else {
				log.Printf("[INFO] Delete DirectMail Dedicated Ip Pool success: %s", name)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageIndex"] = request["PageIndex"].(int) + 1
	}
	return nil
}

// Creating a dedicated IP pool requires purchased dedicated IP instances in
// the account (pool quota derives from purchased IPs); the shared tf
// acceptance test account owns none, so this case skips at runtime unless
// ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN is set.
func TestAccAliCloudDirectMailDedicatedIpPool_basic(t *testing.T) {
	if os.Getenv("ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN") == "" {
		t.Skip("requires purchased dedicated IPs; set ALICLOUD_DIRECT_MAIL_DEDICATED_IP_RUN=1 to run")
	}
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudDirectMailDedicatedIpPoolMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailDedicatedIpPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-direct-mail-pool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDirectMailDedicatedIpPoolBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": "ip-1,ip-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"ip_count": "0",
						"ips.#":    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"buy_resource_ids": "ip-2,ip-3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"buy_resource_ids": "ip-2,ip-3",
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

func TestAccAliCloudDirectMailDedicatedIpPoolsDataSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_dedicated_ip_pool.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DirectMailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDirectMailDedicatedIpPool")
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_direct_mail_dedicated_ip_pools.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_direct_mail_dedicated_ip_pools.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_direct_mail_dedicated_ip_pools.default", "pools.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDirectMailDedicatedIpPoolsDataSourceConfig = `
data "alicloud_direct_mail_dedicated_ip_pools" "default" {
}
`

var AlicloudDirectMailDedicatedIpPoolMap0 = map[string]string{
	"create_time": CHECKSET,
	"ip_count":    CHECKSET,
}

func AlicloudDirectMailDedicatedIpPoolBasicDependence0(name string) string {
	return ""
}
