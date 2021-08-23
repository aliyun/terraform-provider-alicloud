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
	resource.AddTestSweepers("alicloud_vswitch", &resource.Sweeper{
		Name: "alicloud_vswitch",
		F:    testSweepVSwitches,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_db_instance",
			"alicloud_slb_load_balancer",
			"alicloud_ess_scalinggroup",
			"alicloud_fc_service",
			"alicloud_cs_kubernetes",
			"alicloud_kvstore_instance",
			"alicloud_route_table_attachment",
			//"alicloud_havip",
			"alicloud_ecs_network_interface",
			"alicloud_drds_instance",
			"alicloud_elasticsearch_instance",
			"alicloud_vpn_gateway",
			"alicloud_mongodb_instance",
			"alicloud_mongodb_sharding_instance",
			"alicloud_gpdb_instance",
			"alicloud_yundun_bastionhost_instance",
			"alicloud_yundun_dbaudit_instance",
			"alicloud_emr_cluster",
			"alicloud_polardb_cluster",
			"alicloud_hbase_instance",
			"alicloud_cassandra_cluster",
			"alicloud_network_acl",
		},
	})
}

func testSweepVSwitches(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeVSwitches"
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId
	vswitches := make([]map[string]interface{}, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VSwitch in service list: %s", err)
			return nil
		}

		resp, err := jsonpath.Get("$.VSwitches.VSwitch", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VSwitches.VSwitch", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			vswitches = append(vswitches, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	sweeped := false
	service := VpcService{client}
	for _, vsw := range vswitches {
		name := fmt.Sprint(vsw["VSwitchName"])
		id := fmt.Sprint(vsw["VSwitchId"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a vswitch name is set by other service, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(fmt.Sprint(vsw["VpcId"]), ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VSwitch: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting VSwitch: %s (%s)", name, id)
		if err := service.sweepVSwitch(id); err != nil {
			log.Printf("[ERROR] Failed to delete VSwitch (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudVSwitch_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vswitch.default"
	ra := resourceAttrInit(resourceId, AlicloudVswitchMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVswitch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svswitch%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVswitchBasicDependence0)
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
					"zone_id":    "${data.alicloud_zones.default.zones.0.id}",
					"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
					"cidr_block": "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 4, 2)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":    CHECKSET,
						"vpc_id":     CHECKSET,
						"cidr_block": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
					"vswitch_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  name,
						"vswitch_name": name,
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
		},
	})
}

var AlicloudVswitchMap0 = map[string]string{}

func AlicloudVswitchBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
`, name)
}
