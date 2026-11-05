package selectdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/selectdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_selectdb_db_instance", &resource.Sweeper{
		Name: "alicloud_selectdb_db_instance",
		F:    testSweepSelectDBDbInstance,
	})
}

func testSweepSelectDBDbInstance(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	selectDBService := selectdb.NewService(client)
	instanceResp, err := selectDBService.DescribeDbInstances("", nil)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_selectdb_db_instances", errs.AlibabaCloudSdkGoERROR)
	}

	var objects []map[string]interface{}

	for _, item := range instanceResp {
		name := item["Description"].(string)
		skip := true
		if !sweep.SweepAll() {
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
		_, err := selectDBService.DeleteInstance(id["DBInstanceId"].(string))
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_selectdb_db_instances", errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func TestAccAliCloudSelectDBDbInstance_basic_info_upgrade_major_version(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSelectDBDbInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return selectdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSelectDBDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class":       "selectdb.g.2xlarge",
					"db_instance_description": name,
					"cache_size":              "200",
					"payment_type":            "PayAsYouGo",
					"engine_minor_version":    "4.0",
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
						"db_instance_class":           "selectdb.g.2xlarge",
						"cache_size":                  "200",
						"payment_type":                "PayAsYouGo",
						"engine_minor_version":        "4.0",
						"zone_id":                     acctest.CHECKSET,
						"vpc_id":                      acctest.CHECKSET,
						"vswitch_id":                  acctest.CHECKSET,
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
					"engine_minor_version": "4.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_minor_version": "4.1",
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
					"enable_public_network", "admin_pass"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDbInstance_basic_payment_modify_upgrade_minor_version(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSelectDBDbInstanceMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return selectdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSelectDBDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class":       "selectdb.g.2xlarge",
					"db_instance_description": name,
					"cache_size":              "200",
					"engine_minor_version":    "4.0",
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
						"db_instance_class": "selectdb.g.2xlarge",
						"cache_size":        "200",
						"payment_type":      "Subscription",
						"period":            "Month",
						"period_time":       "1",
						"zone_id":           acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
						"vswitch_id":        acctest.CHECKSET,
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
					"engine_minor_version": "4.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_minor_version": "4.1",
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
				ImportStateVerifyIgnore: []string{"gmt_modified", "cache_size", "db_instance_class",
					"enable_public_network", "period", "period_time", "instance_net_infos", "admin_pass"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDbInstance_deprecated_upgrade_minor_version(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSelectDBDbInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return selectdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSelectDBDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class":             "selectdb.g.2xlarge",
					"db_instance_description":       name,
					"cache_size":                    "200",
					"upgraded_engine_minor_version": "4.0",
					"admin_pass":                    "test_123",
					"payment_type":                  "PayAsYouGo",
					"zone_id":                       "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"vpc_id":                        "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"vswitch_id":                    "${data.alicloud_vswitches.default.vswitches.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":             "selectdb.g.2xlarge",
						"cache_size":                    "200",
						"upgraded_engine_minor_version": "4.0",
						"payment_type":                  "PayAsYouGo",
						"zone_id":                       acctest.CHECKSET,
						"vpc_id":                        acctest.CHECKSET,
						"vswitch_id":                    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgraded_engine_minor_version": "4.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upgraded_engine_minor_version": "4.1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"gmt_modified", "cache_size", "db_instance_class",
					"enable_public_network", "instance_net_infos", "admin_pass"},
			},
		},
	})
}

var AliCloudSelectDBDbInstanceMap0 = map[string]string{
	"db_instance_class": acctest.CHECKSET,
	"cache_size":        acctest.CHECKSET,
	"payment_type":      acctest.CHECKSET,
	"zone_id":           acctest.CHECKSET,
	"vpc_id":            acctest.CHECKSET,
	"vswitch_id":        acctest.CHECKSET,
}

var AliCloudSelectDBDbInstanceMap1 = map[string]string{
	"db_instance_class": acctest.CHECKSET,
	"cache_size":        acctest.CHECKSET,
	"payment_type":      acctest.CHECKSET,
	"period":            acctest.CHECKSET,
	"period_time":       acctest.CHECKSET,
	"zone_id":           acctest.CHECKSET,
	"vpc_id":            acctest.CHECKSET,
	"vswitch_id":        acctest.CHECKSET,
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
