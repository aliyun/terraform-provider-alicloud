package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_pvtz_endpoint",
		&resource.Sweeper{
			Name: "alicloud_pvtz_endpoint",
			F:    testSweepPrivateZoneEndpoint,
		})
}

func testSweepPrivateZoneEndpoint(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeResolverEndpoints"
	request := map[string]interface{}{}

	request["Lang"] = "en"
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Endpoints", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Endpoints", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["Name"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping PrivateZone Endpoint: %s", item["Name"].(string))
				continue
			}
			action := "DeleteResolverEndpoint"
			request := map[string]interface{}{
				"EndpointId": item["Id"],
			}
			request["ClientToken"] = buildClientToken(action)
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				log.Printf("[ERROR] Failed to delete PrivateZone Endpoint (%s): %s", item["Id"].(string), err)
			}
			log.Printf("[INFO] Delete PrivateZone Endpoint success: %s ", item["Id"].(string))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudPrivateZoneEndpoint_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateZoneEndpointMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateZoneEndpointBasicDependence0)
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
					"vpc_id":            "${alicloud_vpc.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"vpc_region_id":     defaultRegionToTest,
					"endpoint_name":     name,
					"ip_configs": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.default[0].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[0].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[0].id}",
						},
						{
							"zone_id":    "${alicloud_vswitch.default[1].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[1].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[1].id}",
						},
						{
							"zone_id":    "${alicloud_vswitch.default[2].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[2].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[2].id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"security_group_id": CHECKSET,
						"vpc_region_id":     CHECKSET,
						"endpoint_name":     name,
						"ip_configs.#":      "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_configs": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.default[2].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[2].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[2].id}",
						},
						{
							"zone_id":    "${alicloud_vswitch.default[3].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[3].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[3].id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_configs.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudPrivateZoneEndpointMap0 = map[string]string{}

func AlicloudPrivateZoneEndpointBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_security_group" "default" {
  vpc_id      = alicloud_vpc.default.id
  name        = var.name
}

resource "alicloud_vswitch" "default" {
  count      = 4
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}
`, name)
}
