package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEnsInstance_basic3529(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap3529)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence3529)
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
					"schedule_area_level":        "Region",
					"instance_type":              "ens.sn1.tiny",
					"internet_max_bandwidth_out": "100",
					"payment_type":               "Subscription",
					"instance_name":              name,
					"ens_region_id":              "cn-wuxi-telecom_unicom_cmcc-2",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"image_id": "centos_6_08_64_20G_alibase_20171208",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_area_level":        "Region",
						"instance_type":              "ens.sn1.tiny",
						"internet_max_bandwidth_out": "100",
						"payment_type":               "Subscription",
						"instance_name":              name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "1",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
						{
							"size":     "30",
							"category": "cloud_efficiency",
						},
						{
							"size":     "40",
							"category": "cloud_efficiency",
						},
					},
					"public_ip_identification":   "true",
					"period_unit":                "Month",
					"scheduling_strategy":        "Concentrate",
					"schedule_area_level":        "Region",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"carrier":                    "cmcc",
					"instance_type":              "ens.sn1.tiny",
					"host_name":                  "testHost80",
					"password":                   "Test123456@@",
					"net_district_code":          "100102",
					"internet_charge_type":       "95BandwidthByMonth",
					"instance_name":              name + "_update",
					"internet_max_bandwidth_out": "100",
					"ens_region_id":              "cn-wuxi-telecom_unicom_cmcc-2",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_price_strategy": "PriceHighPriority",
					"user_data":                 "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"instance_charge_strategy":  "user",
					"payment_type":              "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":                     "1",
						"data_disk.#":                "3",
						"public_ip_identification":   "true",
						"period_unit":                "Month",
						"scheduling_strategy":        "Concentrate",
						"schedule_area_level":        "Region",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"carrier":                    "cmcc",
						"instance_type":              "ens.sn1.tiny",
						"host_name":                  "testHost80",
						"password":                   "Test123456@@",
						"net_district_code":          "100102",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name + "_update",
						"internet_max_bandwidth_out": "100",
						"ens_region_id":              "cn-wuxi-telecom_unicom_cmcc-2",
						"scheduling_price_strategy":  "PriceHighPriority",
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"instance_charge_strategy":   "user",
						"payment_type":               "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "carrier", "instance_charge_strategy", "internet_charge_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data", "amount"},
			},
		},
	})
}

var AlicloudEnsInstanceMap3529 = map[string]string{
	"host_name": CHECKSET,
	"status":    CHECKSET,
}

func AlicloudEnsInstanceBasicDependence3529(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3497
func TestAccAliCloudEnsInstance_basic3497(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap3497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence3497)
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
					"schedule_area_level":        "Region",
					"instance_type":              "ens.sn1.tiny",
					"internet_max_bandwidth_out": "100",
					"payment_type":               "PayAsYouGo",
					"instance_name":              name,
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"ens_region_id": "cn-hefei-cmcc-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_area_level":        "Region",
						"instance_type":              "ens.sn1.tiny",
						"internet_max_bandwidth_out": "100",
						"payment_type":               "PayAsYouGo",
						"instance_name":              name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "1",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "local_ssd",
						},
					},
					"public_ip_identification":   "true",
					"period_unit":                "Month",
					"scheduling_strategy":        "Concentrate",
					"schedule_area_level":        "Region",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"carrier":                    "cmcc",
					"instance_type":              "ens.sn1.tiny",
					"host_name":                  "testHost72",
					"password":                   "Test123456@@",
					"net_district_code":          "100102",
					"internet_charge_type":       "95BandwidthByMonth",
					"instance_name":              name + "_update",
					"internet_max_bandwidth_out": "100",
					"ens_region_id":              "cn-hefei-cmcc-2",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_price_strategy": "PriceHighPriority",
					"user_data":                 "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"instance_charge_strategy":  "user",
					"payment_type":              "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":                     "1",
						"data_disk.#":                "1",
						"public_ip_identification":   "true",
						"period_unit":                "Month",
						"scheduling_strategy":        "Concentrate",
						"schedule_area_level":        "Region",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"carrier":                    "cmcc",
						"instance_type":              "ens.sn1.tiny",
						"host_name":                  "testHost72",
						"password":                   "Test123456@@",
						"net_district_code":          "100102",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name + "_update",
						"internet_max_bandwidth_out": "100",
						"ens_region_id":              "cn-hefei-cmcc-2",
						"scheduling_price_strategy":  "PriceHighPriority",
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"instance_charge_strategy":   "user",
						"payment_type":               "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "carrier", "instance_charge_strategy", "internet_charge_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data", "amount"},
			},
		},
	})
}

var AlicloudEnsInstanceMap3497 = map[string]string{
	"host_name": CHECKSET,
	"status":    CHECKSET,
}

func AlicloudEnsInstanceBasicDependence3497(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3529  twin
func TestAccAliCloudEnsInstance_basic3529_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap3529)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence3529)
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
					"period": "1",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
						{
							"size":     "30",
							"category": "cloud_efficiency",
						},
						{
							"size":     "40",
							"category": "cloud_efficiency",
						},
					},
					"public_ip_identification":   "true",
					"period_unit":                "Month",
					"scheduling_strategy":        "Concentrate",
					"schedule_area_level":        "Region",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"carrier":                    "cmcc",
					"instance_type":              "ens.sn1.tiny",
					"host_name":                  "testHost80",
					"password":                   "Test123456@@",
					"net_district_code":          "100102",
					"internet_charge_type":       "95BandwidthByMonth",
					"instance_name":              name,
					"internet_max_bandwidth_out": "100",
					"ens_region_id":              "cn-wuxi-telecom_unicom_cmcc-2",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_price_strategy": "PriceHighPriority",
					"user_data":                 "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"instance_charge_strategy":  "user",
					"payment_type":              "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":                     "1",
						"data_disk.#":                "3",
						"public_ip_identification":   "true",
						"period_unit":                "Month",
						"scheduling_strategy":        "Concentrate",
						"schedule_area_level":        "Region",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"carrier":                    "cmcc",
						"instance_type":              "ens.sn1.tiny",
						"host_name":                  "testHost80",
						"password":                   "Test123456@@",
						"net_district_code":          "100102",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"internet_max_bandwidth_out": "100",
						"ens_region_id":              "cn-wuxi-telecom_unicom_cmcc-2",
						"scheduling_price_strategy":  "PriceHighPriority",
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"instance_charge_strategy":   "user",
						"payment_type":               "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "carrier", "instance_charge_strategy", "internet_charge_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data", "amount"},
			},
		},
	})
}

// Case 3497  twin
func TestAccAliCloudEnsInstance_basic3497_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap3497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence3497)
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
					"period": "1",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "local_ssd",
						},
					},
					"public_ip_identification":   "true",
					"period_unit":                "Month",
					"scheduling_strategy":        "Concentrate",
					"schedule_area_level":        "Region",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"carrier":                    "cmcc",
					"instance_type":              "ens.sn1.tiny",
					"host_name":                  "testHost72",
					"password":                   "Test123456@@",
					"net_district_code":          "100102",
					"internet_charge_type":       "95BandwidthByMonth",
					"instance_name":              name,
					"internet_max_bandwidth_out": "100",
					"ens_region_id":              "cn-hefei-cmcc-2",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_price_strategy": "PriceHighPriority",
					"user_data":                 "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"instance_charge_strategy":  "user",
					"payment_type":              "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":                     "1",
						"data_disk.#":                "1",
						"public_ip_identification":   "true",
						"period_unit":                "Month",
						"scheduling_strategy":        "Concentrate",
						"schedule_area_level":        "Region",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"carrier":                    "cmcc",
						"instance_type":              "ens.sn1.tiny",
						"host_name":                  "testHost72",
						"password":                   "Test123456@@",
						"net_district_code":          "100102",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"internet_max_bandwidth_out": "100",
						"ens_region_id":              "cn-hefei-cmcc-2",
						"scheduling_price_strategy":  "PriceHighPriority",
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"instance_charge_strategy":   "user",
						"payment_type":               "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "carrier", "instance_charge_strategy", "internet_charge_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data", "amount"},
			},
		},
	})
}

// Case 实例创建_后付费区域调度 5654
func TestAccAliCloudEnsInstance_basic5654(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap5654)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence5654)
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
					"net_district_code":    "100102",
					"amount":               "1",
					"period":               "1",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name,
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_strategy":        "Disperse",
					"schedule_area_level":        "Big",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceLowPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "PayAsYouGo",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"ip_type":                    "ipv4",
					"auto_renew":                 "false",
					"password_inherit":           "true",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"carrier":                    "cmcc",
					"billing_cycle":              "Day",
					"instance_charge_strategy":   "instance",
					"period_unit":                "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"amount":                     "1",
						"schedule_area_level":        "Big",
						"internet_max_bandwidth_out": "10",
						"payment_type":               "PayAsYouGo",
						"instance_type":              "ens.sn1.stiny",
						"instance_name":              name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "InstanceHostName_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "InstanceHostName_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "ens.sn1.tiny",
					"status":        "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "ens.sn1.tiny",
						"status":        "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"net_district_code":    "100102",
					"amount":               "1",
					"period":               "1",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name + "_update",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_strategy":        "Disperse",
					"schedule_area_level":        "Big",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceLowPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "PayAsYouGo",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"ip_type":                    "ipv4",
					"auto_renew":                 "false",
					"password_inherit":           "true",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"carrier":                    "cmcc",
					"billing_cycle":              "Day",
					"instance_charge_strategy":   "instance",
					"period_unit":                "Month",
					"status":                     "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"net_district_code":          "100102",
						"amount":                     "1",
						"period":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name + "_update",
						"scheduling_strategy":        "Disperse",
						"schedule_area_level":        "Big",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceLowPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "PayAsYouGo",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"ip_type":                    "ipv4",
						"auto_renew":                 "false",
						"password_inherit":           "true",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"carrier":                    "cmcc",
						"billing_cycle":              "Day",
						"instance_charge_strategy":   "instance",
						"period_unit":                "Month",
						"status":                     "Running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "auto_use_coupon", "billing_cycle", "carrier", "force_stop", "include_data_disks", "instance_charge_strategy", "internet_charge_type", "ip_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

var AlicloudEnsInstanceMap5654 = map[string]string{
	"host_name":     CHECKSET,
	"status":        CHECKSET,
	"instance_name": CHECKSET,
}

func AlicloudEnsInstanceBasicDependence5654(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 实例创建_后付费区域调度 5654  twin
func TestAccAliCloudEnsInstance_basic5654_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap5654)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence5654)
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
					"net_district_code":    "100102",
					"amount":               "1",
					"period":               "1",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name,
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_strategy":        "Disperse",
					"schedule_area_level":        "Big",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceLowPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "PayAsYouGo",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"ip_type":                    "ipv4",
					"auto_renew":                 "false",
					"password_inherit":           "true",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"carrier":                    "cmcc",
					"billing_cycle":              "Day",
					"instance_charge_strategy":   "instance",
					"period_unit":                "Month",
					"status":                     "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"net_district_code":          "100102",
						"amount":                     "1",
						"period":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"scheduling_strategy":        "Disperse",
						"schedule_area_level":        "Big",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceLowPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "PayAsYouGo",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"ip_type":                    "ipv4",
						"auto_renew":                 "false",
						"password_inherit":           "true",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"carrier":                    "cmcc",
						"billing_cycle":              "Day",
						"instance_charge_strategy":   "instance",
						"period_unit":                "Month",
						"status":                     "Stopped",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "auto_use_coupon", "billing_cycle", "carrier", "force_stop", "include_data_disks", "instance_charge_strategy", "internet_charge_type", "ip_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

// Test Ens Instance. >>> Resource test cases, automatically generated.
// Case 实例创建_预付费_网络参数 5657
func TestAccAliCloudEnsInstance_basic5657(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap5657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence5657)
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
					"user_data":            "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"amount":               "1",
					"period":               "1",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name,
					"system_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"scheduling_strategy": "Concentrate",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"schedule_area_level":        "Region",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceHighPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "Subscription",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"period_unit":                "Month",
					"ip_type":                    "ipv4",
					"auto_renew":                 "false",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"ens_region_id":              "cn-nanjing-cmcc",
					"password_inherit":           "true",
					"private_ip_address":         "192.168.2.4",
					"net_work_id":                "${alicloud_ens_network.创建网络.id}",
					"vswitch_id":                 "${alicloud_ens_vswitch.创建交换机.id}",
					"security_id":                "${alicloud_ens_security_group.创建安全组.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"amount":                     "1",
						"period":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"scheduling_strategy":        "Concentrate",
						"data_disk.#":                "1",
						"schedule_area_level":        "Region",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceHighPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "Subscription",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"period_unit":                "Month",
						"ip_type":                    "ipv4",
						"auto_renew":                 "false",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"ens_region_id":              "cn-nanjing-cmcc",
						"password_inherit":           "true",
						"private_ip_address":         "192.168.2.4",
						"net_work_id":                CHECKSET,
						"vswitch_id":                 CHECKSET,
						"security_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "InstanceHostName_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "InstanceHostName_autotest",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "auto_use_coupon", "billing_cycle", "carrier", "force_stop", "include_data_disks", "instance_charge_strategy", "internet_charge_type", "ip_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

var AlicloudEnsInstanceMap5657 = map[string]string{
	"host_name":     CHECKSET,
	"status":        CHECKSET,
	"instance_name": CHECKSET,
}

func AlicloudEnsInstanceBasicDependence5657(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_network" "创建网络" {
  network_name = var.name

  description   = "NetworkDescription_autotest"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-nanjing-cmcc"
}

resource "alicloud_ens_vswitch" "创建交换机" {
  description  = "VSwitchDescription_autotest"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = var.name

  ens_region_id = "cn-nanjing-cmcc"
  network_id    = alicloud_ens_network.创建网络.id
}

resource "alicloud_ens_security_group" "创建安全组" {
  description         = "SecurityGroupDescription_autotest"
  security_group_name = var.name

}


`, name)
}

// Case 实例创建 5608
func TestAccAliCloudEnsInstance_basic5608(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap5608)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence5608)
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
					"user_data":            "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"amount":               "1",
					"period":               "2",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name,
					"system_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"scheduling_strategy": "Concentrate",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"schedule_area_level":        "Region",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceHighPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "Subscription",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"password":                   "12345678abcABC",
					"period_unit":                "Month",
					"ip_type":                    "ipv4",
					"auto_renew":                 "true",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"ens_region_id":              "cn-nanjing-cmcc",
					"password_inherit":           "false",
					"status":                     "Stopped",
					"force_stop":                 "true",
					"include_data_disks":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"amount":                     "1",
						"period":                     "2",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"scheduling_strategy":        "Concentrate",
						"data_disk.#":                "1",
						"schedule_area_level":        "Region",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceHighPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "Subscription",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"password":                   "12345678abcABC",
						"period_unit":                "Month",
						"ip_type":                    "ipv4",
						"auto_renew":                 "true",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"ens_region_id":              "cn-nanjing-cmcc",
						"password_inherit":           "false",
						"status":                     "Stopped",
						"force_stop":                 "true",
						"include_data_disks":         "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "InstanceHostName_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "InstanceHostName_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "12345678abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "12345678abcABC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
					"status":        "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
						"status":        "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data":            "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"amount":               "1",
					"period":               "1",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name + "_update",
					"system_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"scheduling_strategy": "Concentrate",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"schedule_area_level":        "Region",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceHighPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "Subscription",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"password":                   "12345678abcABC",
					"period_unit":                "Month",
					"ip_type":                    "ipv4",
					"auto_renew":                 "false",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"ens_region_id":              "cn-nanjing-cmcc",
					"password_inherit":           "false",
					"status":                     "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"amount":                     "1",
						"period":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name + "_update",
						"scheduling_strategy":        "Concentrate",
						"data_disk.#":                "1",
						"schedule_area_level":        "Region",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceHighPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "Subscription",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"password":                   "12345678abcABC",
						"period_unit":                "Month",
						"ip_type":                    "ipv4",
						"auto_renew":                 "false",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"ens_region_id":              "cn-nanjing-cmcc",
						"password_inherit":           "false",
						"status":                     "Running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "auto_use_coupon", "billing_cycle", "carrier", "force_stop", "include_data_disks", "instance_charge_strategy", "internet_charge_type", "ip_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

var AlicloudEnsInstanceMap5608 = map[string]string{
	"host_name":     CHECKSET,
	"status":        CHECKSET,
	"instance_name": CHECKSET,
}

func AlicloudEnsInstanceBasicDependence5608(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 实例创建_预付费_网络参数 5657  twin
func TestAccAliCloudEnsInstance_basic5657_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap5657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence5657)
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
					"user_data":            "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"amount":               "1",
					"period":               "1",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name,
					"system_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"scheduling_strategy": "Concentrate",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"schedule_area_level":        "Region",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceHighPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "Subscription",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"period_unit":                "Month",
					"ip_type":                    "ipv4",
					"auto_renew":                 "false",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"ens_region_id":              "cn-nanjing-cmcc",
					"password_inherit":           "true",
					"private_ip_address":         "192.168.2.4",
					"net_work_id":                "${alicloud_ens_network.创建网络.id}",
					"vswitch_id":                 "${alicloud_ens_vswitch.创建交换机.id}",
					"security_id":                "${alicloud_ens_security_group.创建安全组.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"amount":                     "1",
						"period":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"scheduling_strategy":        "Concentrate",
						"data_disk.#":                "1",
						"schedule_area_level":        "Region",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceHighPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "Subscription",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"period_unit":                "Month",
						"ip_type":                    "ipv4",
						"auto_renew":                 "false",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"ens_region_id":              "cn-nanjing-cmcc",
						"password_inherit":           "true",
						"private_ip_address":         "192.168.2.4",
						"net_work_id":                CHECKSET,
						"vswitch_id":                 CHECKSET,
						"security_id":                CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "auto_use_coupon", "billing_cycle", "carrier", "force_stop", "include_data_disks", "instance_charge_strategy", "internet_charge_type", "ip_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

// Case 实例创建 5608  twin
func TestAccAliCloudEnsInstance_basic5608_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap5608)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence5608)
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
					"user_data":            "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"amount":               "1",
					"period":               "2",
					"internet_charge_type": "95BandwidthByMonth",
					"instance_name":        name,
					"system_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"scheduling_strategy": "Concentrate",
					"data_disk": []map[string]interface{}{
						{
							"size":     "20",
							"category": "cloud_efficiency",
						},
					},
					"schedule_area_level":        "Region",
					"internet_max_bandwidth_out": "10",
					"public_ip_identification":   "true",
					"scheduling_price_strategy":  "PriceHighPriority",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"payment_type":               "Subscription",
					"instance_type":              "ens.sn1.stiny",
					"host_name":                  "InstanceHostName_autotest",
					"password":                   "12345678abcABC",
					"period_unit":                "Month",
					"ip_type":                    "ipv4",
					"auto_renew":                 "true",
					"unique_suffix":              "false",
					"auto_use_coupon":            "true",
					"ens_region_id":              "cn-nanjing-cmcc",
					"password_inherit":           "false",
					"status":                     "Stopped",
					"force_stop":                 "true",
					"include_data_disks":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"amount":                     "1",
						"period":                     "2",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"scheduling_strategy":        "Concentrate",
						"data_disk.#":                "1",
						"schedule_area_level":        "Region",
						"internet_max_bandwidth_out": "10",
						"public_ip_identification":   "true",
						"scheduling_price_strategy":  "PriceHighPriority",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"payment_type":               "Subscription",
						"instance_type":              "ens.sn1.stiny",
						"host_name":                  "InstanceHostName_autotest",
						"password":                   "12345678abcABC",
						"period_unit":                "Month",
						"ip_type":                    "ipv4",
						"auto_renew":                 "true",
						"unique_suffix":              "false",
						"auto_use_coupon":            "true",
						"ens_region_id":              "cn-nanjing-cmcc",
						"password_inherit":           "false",
						"status":                     "Stopped",
						"force_stop":                 "true",
						"include_data_disks":         "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "auto_use_coupon", "billing_cycle", "carrier", "force_stop", "include_data_disks", "instance_charge_strategy", "internet_charge_type", "ip_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

// Test Ens Instance. <<< Resource test cases, automatically generated.
