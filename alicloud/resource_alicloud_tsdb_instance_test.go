package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_tsdb_instance", &resource.Sweeper{
		Name: "alicloud_tsdb_instance",
		F:    testSweepTsdbInstance,
	})
}

func testSweepTsdbInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := make(map[string]interface{}, 0)
	action := "DescribeHiTSDBInstanceList"
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	var instances []string
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve tsdb instance service list: %s", err)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if v, ok := item["InstanceAlias"]; !ok || v.(string) == "" {
				continue
			}
			instances = append(instances, fmt.Sprint(item["InstanceId"], ":", item["InstanceAlias"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, instance := range instances {
		instanceId := strings.Split(instance, ":")[0]
		instanceName := strings.Split(instance, ":")[1]
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping tsdb instance: %s ", instanceId)
			continue
		}

		action := "DeleteHiTSDBInstance"
		var response map[string]interface{}
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"InstanceId": instanceId,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve tsdb instance %s %v", instanceId, err)
			continue
		}
		log.Printf("[INFO] Delete tsdb instance: %s ", instanceId)
	}
	return nil
}

func TestAccAlicloudTsdbInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tsdb_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudTsdbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTsdbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccTsdbInstance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudTsdbInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.TsdbInstanceSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":     "PayAsYouGo",
					"instance_class":   "tsdb.1x.basic",
					"instance_storage": "50",
					"vswitch_id":       "${local.vswitch_id}",
					"engine_type":      "tsdb_tsdb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":     "PayAsYouGo",
						"instance_class":   "tsdb.1x.basic",
						"instance_storage": "50",
						"vswitch_id":       CHECKSET,
						"engine_type":      "tsdb_tsdb",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": name + "change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "tsdb.4x.basic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "tsdb.4x.basic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":     "PayAsYouGo",
					"instance_class":   "tsdb.1x.basic",
					"instance_storage": "150",
					"vswitch_id":       "${local.vswitch_id}",
					"engine_type":      "tsdb_tsdb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":     "PayAsYouGo",
						"instance_class":   "tsdb.1x.basic",
						"instance_storage": "150",
						"vswitch_id":       CHECKSET,
						"engine_type":      "tsdb_tsdb",
					}),
				),
			},
		},
	})
}

var AlicloudTsdbInstanceMap = map[string]string{
	"status": "ACTIVATION",
}

func AlicloudTsdbInstanceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_tsdb_zones" "default" {}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_tsdb_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_tsdb_zones.default.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

`, name)
}
