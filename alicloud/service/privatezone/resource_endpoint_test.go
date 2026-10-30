package privatezone_test

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	client, err := privatezoneTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeResolverEndpoints"
	request := map[string]interface{}{}

	request["Lang"] = "en"
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Endpoints", response)
		totalCount, formatErr := convert.FormatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
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
			request["ClientToken"] = helper.BuildClientToken(action)
			_, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete PrivateZone Endpoint (%s): %s", item["Id"].(string), err)
			}
			log.Printf("[INFO] Delete PrivateZone Endpoint success: %s ", item["Id"].(string))
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudPvtzEndpoint_basic0(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateZoneEndpointMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEndpoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateZoneEndpointBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip":                acctest.REMOVEKEY,
					"vpc_id":            "${alicloud_vpc.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"vpc_region_id":     acctest.DefaultRegionToTest,
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
						"vpc_id":            acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
						"vpc_region_id":     acctest.CHECKSET,
						"endpoint_name":     name,
						"ip_configs.#":      "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip":            acctest.REMOVEKEY,
					"endpoint_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_name": name + "1",
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
				// The endpoint IP is assigned by the service and cannot be
				// pre-provisioned in this test. Removing the required
				// endpoint_name stops the plan before the placeholder
				// ip_configs entry reaches the API.
				Config: testAccConfig(map[string]interface{}{
					"endpoint_name": acctest.REMOVEKEY,
					"ip_configs": []map[string]interface{}{
						{
							"ip":         "",
							"zone_id":    "${alicloud_vswitch.default[0].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[0].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[0].id}",
						},
					},
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				// The ExpectError step above omits the required endpoint_name,
				// leaving an invalid config as the last written one. Restore a
				// valid Config (matching the Step 4 state) as the final apply
				// step so the post-test destroy runs against a valid plan. This
				// apply is a no-op against the existing resource. ip is
				// Optional+Computed and service-assigned, omitted as in Step 1.
				Config: testAccConfig(map[string]interface{}{
					"endpoint_name": name,
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
						"endpoint_name": name,
					}),
				),
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

// lintignore: R001
