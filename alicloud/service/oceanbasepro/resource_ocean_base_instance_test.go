package oceanbasepro_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oceanbasepro"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOceanBaseInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oceanbasepro.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOceanBaseInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_name":     "${var.name}",
					"series":            "normal",
					"auto_renew":        "false",
					"disk_size":         "100",
					"payment_type":      "PayAsYouGo",
					"instance_class":    "8C32GB",
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"series":         "normal",
						"auto_renew":     "false",
						"disk_size":      "100",
						"instance_class": acctest.CHECKSET,
						//"resource_group_id": CHECKSET,
						"zones.#":      "3",
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_num": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_num": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size":      "200",
					"instance_class": "14C70GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size":      "200",
						"instance_class": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size":      "100",
					"instance_class": "8C32G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size":      "100",
						"instance_class": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

func TestAccAliCloudOceanBaseInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oceanbasepro.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOceanBaseInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  "${var.name}",
					"series":         "normal",
					"auto_renew":     "false",
					"disk_size":      "100",
					"payment_type":   "PayAsYouGo",
					"instance_class": "8C32GB",
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"series":         "normal",
						"auto_renew":     "false",
						"disk_size":      "100",
						"instance_class": acctest.CHECKSET,
						//"resource_group_id": CHECKSET,
						"zones.#":      "3",
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

func TestAccAliCloudOceanBaseInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oceanbasepro.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOceanBaseInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"period":            "1",
					"period_unit":       "Month",
					"auto_renew_period": "1",
					"instance_name":     "${var.name}",
					"series":            "normal",
					"auto_renew":        "false",
					"disk_size":         "300",
					"disk_type":         "cloud_essd_pl1",
					"ob_version":        "4.2.1",
					"payment_type":      "Subscription",
					"instance_class":    "4C16G",
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"series":         "normal",
						"auto_renew":     "false",
						"disk_size":      "300",
						"disk_type":      "cloud_essd_pl1",
						"ob_version":     "4.2.1",
						"instance_class": acctest.CHECKSET,
						//"resource_group_id": CHECKSET,
						"zones.#":      "3",
						"payment_type": "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

var AlicloudOceanBaseInstanceMap0 = map[string]string{
	"status":  acctest.CHECKSET,
	"zones.#": acctest.CHECKSET,
}

func AlicloudOceanBaseInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {}


data "alicloud_resource_manager_resource_groups" "default"{}

`, name)
}

// Test OceanBase Instance. >>> Resource test cases, automatically generated.
// Case tf_0822_qianyu_zhubei 7602
func TestAccAliCloudOceanBaseInstance_basic7602(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap7602)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oceanbasepro.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOceanBaseInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soceanbaseinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOceanBaseInstanceBasicDependence7602)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "4C16G",
					"zones": []string{
						"cn-hangzhou-j", "cn-hangzhou-i"},
					"instance_name":    name,
					"disk_type":        "cloud_essd_pl1",
					"series":           "normal",
					"disk_size":        "100",
					"payment_type":     "PayAsYouGo",
					"ob_version":       "4.2.1",
					"node_num":         "3",
					"cpu_arch":         "X86",
					"primary_region":   "cn-hangzhou",
					"primary_instance": "${alicloud_ocean_base_instance.createZhu.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":   acctest.CHECKSET,
						"zones.#":          "2",
						"instance_name":    name,
						"disk_type":        "cloud_essd_pl1",
						"series":           "normal",
						"disk_size":        "100",
						"payment_type":     "PayAsYouGo",
						"ob_version":       "4.2.1",
						"node_num":         "3",
						"cpu_arch":         "X86",
						"primary_region":   acctest.CHECKSET,
						"primary_instance": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":      "8C32GB",
					"instance_name":       name + "_update",
					"disk_size":           "200",
					"upgrade_spec_native": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":      acctest.CHECKSET,
						"instance_name":       name + "_update",
						"disk_size":           "200",
						"upgrade_spec_native": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"node_num":      "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
						"node_num":      "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  name + "_update",
					"disk_size":      "100",
					"instance_class": "4C16G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name + "_update",
						"disk_size":      "100",
						"instance_class": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "period", "period_unit", "upgrade_spec_native"},
			},
		},
	})
}

var AlicloudOceanBaseInstanceMap7602 = map[string]string{
	"node_num":       acctest.CHECKSET,
	"cpu":            acctest.CHECKSET,
	"disk_type":      acctest.CHECKSET,
	"commodity_code": acctest.CHECKSET,
	"ob_version":     acctest.CHECKSET,
	"status":         acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
	"instance_name":  acctest.CHECKSET,
}

func AlicloudOceanBaseInstanceBasicDependence7602(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ocean_base_instance" "createZhu" {
  instance_class = "4C16G"
  zones          = ["cn-hangzhou-j", "cn-hangzhou-i"]
  instance_name  = var.name
  disk_type      = "cloud_essd_pl1"
  series         = "normal"
  disk_size      = "100"
  payment_type   = "PayAsYouGo"
  ob_version     = "4.2.1"
  period_unit    = "Hour"
  node_num       = "3"
}


`, name)
}

// Test OceanBase Instance. <<< Resource test cases, automatically generated.
