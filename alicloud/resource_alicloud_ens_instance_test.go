package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens Instance. >>> Resource test cases, automatically generated.
// Case 3412
func TestAccAlicloudEnsInstance_basic3412(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap3412)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence3412)
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
					"amount":                     "1",
					"internet_max_bandwidth_out": "100",
					"instance_charge_type":       "PrePaid",
					"instance_name":              name,
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"ens_region_id":              "cn-hefei-cmcc-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_area_level":        "Region",
						"instance_type":              "ens.sn1.tiny",
						"amount":                     "1",
						"internet_max_bandwidth_out": "100",
						"instance_charge_type":       "PrePaid",
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
					"auto_renew":                 "False",
					"scheduling_strategy":        "Concentrate",
					"schedule_area_level":        "Region",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"carrier":                    "cmcc",
					"instance_type":              "ens.sn1.tiny",
					"host_name":                  "testHost56",
					"password":                   "Test123456@@",
					"net_district_code":          "100102",
					"amount":                     "1",
					"internet_charge_type":       "95BandwidthByMonth",
					"instance_name":              name + "_update",
					"internet_max_bandwidth_out": "100",
					"ens_region_id":              "cn-hefei-cmcc-2",
					"instance_charge_type":       "PrePaid",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_price_strategy": "PriceHighPriority",
					"user_data":                 "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"instance_charge_strategy":  "user",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":                     "1",
						"data_disk.#":                "1",
						"public_ip_identification":   "true",
						"period_unit":                "Month",
						"auto_renew":                 "False",
						"scheduling_strategy":        "Concentrate",
						"schedule_area_level":        "Region",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"carrier":                    "cmcc",
						"instance_type":              "ens.sn1.tiny",
						"host_name":                  "testHost56",
						"password":                   "Test123456@@",
						"net_district_code":          "100102",
						"amount":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name + "_update",
						"internet_max_bandwidth_out": "100",
						"ens_region_id":              "cn-hefei-cmcc-2",
						"instance_charge_type":       "PrePaid",
						"scheduling_price_strategy":  "PriceHighPriority",
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"instance_charge_strategy":   "user",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "carrier", "instance_charge_strategy", "instance_charge_type", "internet_charge_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

var AlicloudEnsInstanceMap3412 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
}

func AlicloudEnsInstanceBasicDependence3412(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3412  twin
func TestAccAlicloudEnsInstance_basic3412_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsInstanceMap3412)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsInstanceBasicDependence3412)
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
					"auto_renew":                 "False",
					"scheduling_strategy":        "Concentrate",
					"schedule_area_level":        "Region",
					"image_id":                   "centos_6_08_64_20G_alibase_20171208",
					"carrier":                    "cmcc",
					"instance_type":              "ens.sn1.tiny",
					"host_name":                  "testHost56",
					"password":                   "Test123456@@",
					"net_district_code":          "100102",
					"amount":                     "1",
					"internet_charge_type":       "95BandwidthByMonth",
					"instance_name":              name,
					"internet_max_bandwidth_out": "100",
					"ens_region_id":              "cn-hefei-cmcc-2",
					"instance_charge_type":       "PrePaid",
					"system_disk": []map[string]interface{}{
						{
							"size": "20",
						},
					},
					"scheduling_price_strategy": "PriceHighPriority",
					"user_data":                 "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
					"instance_charge_strategy":  "user",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":                     "1",
						"data_disk.#":                "1",
						"public_ip_identification":   "true",
						"period_unit":                "Month",
						"auto_renew":                 "False",
						"scheduling_strategy":        "Concentrate",
						"schedule_area_level":        "Region",
						"image_id":                   "centos_6_08_64_20G_alibase_20171208",
						"carrier":                    "cmcc",
						"instance_type":              "ens.sn1.tiny",
						"host_name":                  "testHost56",
						"password":                   "Test123456@@",
						"net_district_code":          "100102",
						"amount":                     "1",
						"internet_charge_type":       "95BandwidthByMonth",
						"instance_name":              name,
						"internet_max_bandwidth_out": "100",
						"ens_region_id":              "cn-hefei-cmcc-2",
						"instance_charge_type":       "PrePaid",
						"scheduling_price_strategy":  "PriceHighPriority",
						"user_data":                  "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0",
						"instance_charge_strategy":   "user",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_renew", "carrier", "instance_charge_strategy", "instance_charge_type", "internet_charge_type", "net_district_code", "password", "password_inherit", "period", "period_unit", "public_ip_identification", "schedule_area_level", "scheduling_price_strategy", "scheduling_strategy", "unique_suffix", "user_data"},
			},
		},
	})
}

// Test Ens Instance. <<< Resource test cases, automatically generated.
