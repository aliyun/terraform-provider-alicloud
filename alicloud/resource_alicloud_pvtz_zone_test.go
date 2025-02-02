package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_pvtz_zone", &resource.Sweeper{
		Name: "alicloud_pvtz_zone",
		F:    testSweepPvtzZones,
	})
}

func testSweepPvtzZones(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}
	action := "DescribeZones"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
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
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	for _, v := range zones {
		v := v.(map[string]interface{})
		name := v["ZoneName"].(string)
		id := v["ZoneId"].(string)
		skip := true
		if !sweepAll() {
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
	ra := resourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"user_client_ip":    NOSET,
						"lang":              NOSET,
						"resource_group_id": CHECKSET,
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
					"remark":         REMOVEKEY,
					"proxy_pattern":  "ZONE",
					"lang":           "jp",
					"user_client_ip": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark":         REMOVEKEY,
						"proxy_pattern":  "ZONE",
						"lang":           "jp",
						"user_client_ip": REMOVEKEY,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAliCloudPvtzZone_multi(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone.default.4"
	ra := resourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
	ra := resourceAttrInit(resourceId, pvtzZoneBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzZone")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneConfigDependence)
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
						"resource_group_id": CHECKSET,
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
					"zone_name": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"zone_name": REMOVEKEY,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
	ra := resourceAttrInit(resourceId, pvtzZoneBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzZone")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneConfigDependence)
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
					"zone_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"sync_status":       "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":         name,
						"proxy_pattern":     "ZONE",
						"resource_group_id": CHECKSET,
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
					"user_info":   REMOVEKEY,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
	"zone_name": CHECKSET,
}
