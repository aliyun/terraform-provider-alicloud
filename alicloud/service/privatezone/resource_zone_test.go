package privatezone_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_pvtz_zone", &resource.Sweeper{
		Name: "alicloud_pvtz_zone",
		F:    testSweepPvtzZones,
	})
}

func testSweepPvtzZones(region string) error {
	client, err := privatezoneTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
		"terraform-example",
	}
	action := "DescribeZones"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var zones []interface{}
	for {
		response, err := client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] retrieving Private Zones: %s", err)
		}
		resp, err := jsonpath.Get("$.Zones.Zone", response)
		if err != nil {
			log.Printf("[ERROR] Parsing return parameter error: %s", err)
		}
		result, _ := resp.([]interface{})
		zones = append(zones, result...)
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	for _, v := range zones {
		v := v.(map[string]interface{})
		name := v["ZoneName"].(string)
		id := v["ZoneId"].(string)
		skip := true
		if !privatezoneSweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Private Zone: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Unbinding VPC from Private Zone: %s (%s)", name, id)
		action := "BindZoneVpc"
		request := make(map[string]interface{})
		request["ZoneId"] = id
		vpcs := make([]map[string]interface{}, 0)
		request["Vpcs"] = vpcs
		_, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to unbind VPC from Private Zone (%s (%s)): %s ", name, id, err)
		}

		log.Printf("[INFO] Deleting Private Zone: %s (%s)", name, id)
		action = "DeleteZone"
		request = map[string]interface{}{
			"ZoneId": id,
		}
		_, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Private Zone (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAliCloudPvtzZone_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeZone")

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":         name,
						"proxy_pattern":     "ZONE",
						"user_client_ip":    acctest.NOSET,
						"lang":              acctest.NOSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_client_ip", "lang"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "remark-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "remark-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "remark-test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "remark-test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_pattern": "ZONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_pattern": "ZONE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_client_ip": "172.10.1.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_client_ip": "172.10.1.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_pattern":  "RECORD",
					"lang":           "zh",
					"user_client_ip": "172.10.2.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_pattern":  "RECORD",
						"lang":           "zh",
						"user_client_ip": "172.10.2.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark":         acctest.REMOVEKEY,
					"proxy_pattern":  "ZONE",
					"lang":           "jp",
					"user_client_ip": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark":         acctest.REMOVEKEY,
						"proxy_pattern":  "ZONE",
						"lang":           "jp",
						"user_client_ip": acctest.REMOVEKEY,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAliCloudPvtzZone_dnsGroup(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeZone")

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":         name,
						"dns_group":         acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_client_ip", "lang"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dns_group":         "FAST_ZONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dns_group": "FAST_ZONE",
					}),
				),
			},
		},
	})
}
func TestAccAliCloudPvtzZone_multi(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone.default.4"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeZone")

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_name": fmt.Sprintf("tf-testacc%d${count.index}.test.com", rand),
					"count":     "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudPvtzZone_syncTask(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneBasicMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeZone")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"sync_status":       "OFF",
					"user_info": []map[string]interface{}{
						{
							"user_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.account_id}",
							"region_ids": []string{"cn-beijing", "cn-hangzhou"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":         name,
						"proxy_pattern":     "ZONE",
						"resource_group_id": acctest.CHECKSET,
						"user_info.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sync_status": "OFF",
					"user_info": []map[string]interface{}{
						{
							"user_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.account_id}",
							"region_ids": []string{"cn-beijing", "cn-hangzhou", "cn-chengdu"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sync_status": "OFF",
					"user_info": []map[string]interface{}{
						{
							"user_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.account_id}",
							"region_ids": []string{"cn-beijing", "cn-hangzhou"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":      name,
					"zone_name": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"zone_name": acctest.REMOVEKEY,
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

func TestAccAliCloudPvtzZone_syncTask1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneBasicMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeZone")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"sync_status":       "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":         name,
						"proxy_pattern":     "ZONE",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sync_status": "ON",
					"user_info": []map[string]interface{}{
						{
							"user_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.account_id}",
							"region_ids": []string{"cn-beijing", "cn-hangzhou", "cn-chengdu"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sync_status": "OFF",
					"user_info":   acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sync_status": "OFF",
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
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourcePvtzZoneConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}
`)
}

var pvtzZoneBasicMap = map[string]string{
	"zone_name": acctest.CHECKSET,
}
