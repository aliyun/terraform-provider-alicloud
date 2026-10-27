package mse_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_mse_gateway",
		&resource.Sweeper{
			Name: "alicloud_mse_gateway",
			F:    testSweepMseGateway,
		})
}

func testSweepMseGateway(region string) error {
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.MSEGatewaySupportRegions) {
		log.Printf("[INFO] Skipping Mse Gateway unsupported region: %s", region)
		return nil
	}

	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListGateway"
	request := map[string]interface{}{}

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Data.Result", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Data.Result", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Mse Gateway: %s", item["Name"].(string))
				continue
			}
			action := "DeleteGateway"
			request := map[string]interface{}{
				"GatewayUniqueId": item["GatewayUniqueId"],
			}
			request["DeleteSlb"] = true
			_, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Mse Gateway (%s): %s", item["Name"].(string), err)
			} else {
				log.Printf("[INFO] Delete Mse Gateway success: %s ", item["Name"].(string))
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudMSEGateway_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.MSEGatewaySupportRegions)
	resourceId := "alicloud_mse_gateway.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMSEGatewayMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smsegateway%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMSEGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":              "MSE_GTW_2_4_200_c",
					"replica":           "2",
					"vpc_id":            "${alicloud_vpc.default.id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"gateway_name":      "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":         "MSE_GTW_2_4_200_c",
						"replica":      "2",
						"vpc_id":       acctest.CHECKSET,
						"vswitch_id":   acctest.CHECKSET,
						"gateway_name": name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "data.alicloud_resource_manager_resource_groups.default", "groups.1.id"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": "${var.name}_update",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name + "_update",
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name + "_update",
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "resource_group_id", "data.alicloud_resource_manager_resource_groups.default", "groups.0.id"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enterprise_security_group", "slb_spec", "delete_slb"},
			},
		},
	})
}
func TestAccAliCloudMSEGateway_basic1(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MSEGatewaySupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_mse_gateway.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMSEGatewayMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smsegateway%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMSEGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":                      "MSE_GTW_2_4_200_c",
					"replica":                   "2",
					"vpc_id":                    "${alicloud_vpc.default.id}",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"backup_vswitch_id":         "${alicloud_vswitch.backup.id}",
					"gateway_name":              "${var.name}",
					"enterprise_security_group": "true",
					"slb_spec":                  "slb.s2.small",
					"internet_slb_spec":         "slb.s2.small",
					"delete_slb":                "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":              "MSE_GTW_2_4_200_c",
						"replica":           "2",
						"vpc_id":            acctest.CHECKSET,
						"vswitch_id":        acctest.CHECKSET,
						"backup_vswitch_id": acctest.CHECKSET,
						"gateway_name":      name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enterprise_security_group", "slb_spec", "internet_slb_spec", "delete_slb"},
			},
		},
	})
}

var AlicloudMSEGatewayMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

// AlicloudMSEGatewayBasicDependence0 provisions an isolated VPC with two
// vSwitches (zones.0 for primary, zones.1 for backup) instead of depending on
// the shared default-NODELETING VPC, whose vSwitch availability varies by region
// and zone and caused empty-vSwitch plan failures under CheckoutSupportedRegions.
func AlicloudMSEGatewayBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
	vpc_name   = "${var.name}"
	cidr_block = "172.17.0.0/16"
}
resource "alicloud_vswitch" "default" {
	vswitch_name = "${var.name}"
	cidr_block   = "172.17.3.0/24"
	vpc_id       = alicloud_vpc.default.id
	zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "backup" {
	vswitch_name = "${var.name}-backup"
	cidr_block   = "172.17.4.0/24"
	vpc_id       = alicloud_vpc.default.id
	zone_id      = data.alicloud_zones.default.zones.1.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}
