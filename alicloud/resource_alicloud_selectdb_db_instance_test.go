package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_selectdb_db_instance", &resource.Sweeper{
		Name: "alicloud_selectdb_db_instance",
		F:    testSweepSelectDBDbInstance,
	})
}

func testSweepSelectDBDbInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	selectDBService := SelectDBService{client}
	instanceResp, err := selectDBService.DescribeSelectDBDbInstances("", nil)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instances", AlibabaCloudSdkGoERROR)
	}

	var objects []map[string]interface{}

	for _, item := range instanceResp {
		name := item["Description"].(string)
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(name, prefix) {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}
		objects = append(objects, item)
	}

	for _, id := range objects {
		_, err := selectDBService.DeleteSelectDBInstance(id["DBInstanceId"].(string))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instances", AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func TestAccAliCloudSelectDBDbInstance_basic_info_upgrade_major_version(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class":       "selectdb.xlarge",
					"db_instance_description": name,
					"cache_size":              "200",
					"payment_type":            "PayAsYouGo",
					"zone_id":                 "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"vpc_id":                  "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"vswitch_id":              "${data.alicloud_vswitches.default.vswitches.0.id}",
					"desired_security_ip_lists": []map[string]interface{}{
						{
							"group_name":       "test1",
							"security_ip_list": "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":           "selectdb.xlarge",
						"cache_size":                  "200",
						"payment_type":                "PayAsYouGo",
						"zone_id":                     CHECKSET,
						"vpc_id":                      CHECKSET,
						"vswitch_id":                  CHECKSET,
						"desired_security_ip_lists.#": "1",
						"security_ip_lists.#":         "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_updateAll",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_updateAll",
						"tags.%":                  "2",
						"tags.Created":            "TF-update",
						"tags.For":                "test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_security_ip_lists": []map[string]interface{}{
						{
							"group_name":       "test2",
							"security_ip_list": "192.168.0.3",
						},
						{
							"group_name":       "test1",
							"security_ip_list": "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_security_ip_lists.#": "2",
						"security_ip_lists.#":         "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_pass": "test_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"admin_pass": "test_123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgraded_engine_minor_version": "4.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_minor_version":          "4.0.4",
						"upgraded_engine_minor_version": "4.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cache_size": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cache_size": "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_network": "true",
						"instance_net_infos.#":  "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_network": "false",
						"instance_net_infos.#":  "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"security_ip_lists", "gmt_modified",
					"instance_net_infos", "desired_security_ip_lists", "db_instance_class",
					"enable_public_network", "upgraded_engine_minor_version", "admin_pass"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDbInstance_basic_payment_modify_upgrade_minor_version(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDbInstanceMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class":       "selectdb.xlarge",
					"db_instance_description": name,
					"cache_size":              "200",
					"admin_pass":              "test_123",
					"payment_type":            "Subscription",
					"period":                  "Month",
					"period_time":             "1",
					"zone_id":                 "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"vpc_id":                  "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"vswitch_id":              "${data.alicloud_vswitches.default.vswitches.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.xlarge",
						"cache_size":        "200",
						"payment_type":      "Subscription",
						"period":            "Month",
						"period_time":       "1",
						"zone_id":           CHECKSET,
						"vpc_id":            CHECKSET,
						"vswitch_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgraded_engine_minor_version": "4.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_minor_version":          "4.0.4",
						"upgraded_engine_minor_version": "4.0.4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"gmt_modified", "cache_size", "db_instance_class", "upgraded_engine_minor_version",
					"enable_public_network", "period", "period_time", "instance_net_infos", "engine_minor_version", "admin_pass"},
			},
		},
	})
}

var AliCloudSelectDBDbInstanceMap0 = map[string]string{
	"db_instance_class": CHECKSET,
	"cache_size":        CHECKSET,
	"payment_type":      CHECKSET,
	"zone_id":           CHECKSET,
	"vpc_id":            CHECKSET,
	"vswitch_id":        CHECKSET,
}

var AliCloudSelectDBDbInstanceMap1 = map[string]string{
	"db_instance_class": CHECKSET,
	"cache_size":        CHECKSET,
	"payment_type":      CHECKSET,
	"period":            CHECKSET,
	"period_time":       CHECKSET,
	"zone_id":           CHECKSET,
	"vpc_id":            CHECKSET,
	"vswitch_id":        CHECKSET,
}

func AliCloudSelectDBDbInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}
`, name)
}
