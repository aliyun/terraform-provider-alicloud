package cloudfirewall_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	cloudFirewallInstanceLogisticsConfig = `{\"cityCode\":\"330100\",\"cityName\":\"Hangzhou\",\"contactName\":\"TerraformAcc\",\"countryCode\":\"CN\",\"detailAddress\":\"No.1 Test Road\",\"districtCode\":\"330106\",\"districtName\":\"Xihu District\",\"email\":\"tfacc@example.com\",\"mobilePhone\":\"153564848844\",\"phone\":\"1234567\",\"provCode\":\"330000\",\"provName\":\"Zhejiang\",\"streetCode\":\"33010610\",\"streetName\":\"Zhuantang\",\"zipCode\":\"0000\"}`
	cloudFirewallInstanceLogisticsState  = `{"cityCode":"330100","cityName":"Hangzhou","contactName":"TerraformAcc","countryCode":"CN","detailAddress":"No.1 Test Road","districtCode":"330106","districtName":"Xihu District","email":"tfacc@example.com","mobilePhone":"153564848844","phone":"1234567","provCode":"330000","provName":"Zhejiang","streetCode":"33010610","streetName":"Zhuantang","zipCode":"0000"}`
)

func TestAccAliCloudCloudFirewallInstance_basic0(t *testing.T) {
	ensureCloudFirewallInstanceTestRegion()
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "QueryAvailableInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccCloudFirewallInstancePreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":          "PayAsYouGo",
					"auto_asset_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "PayAsYouGo",
						"auto_asset_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":               "true",
					"modify_type":           "Upgrade",
					"auto_asset_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log":               "true",
						"modify_type":           "Upgrade",
						"auto_asset_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstance_basic0_twin(t *testing.T) {
	ensureCloudFirewallInstanceTestRegion()
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "QueryAvailableInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccCloudFirewallInstancePreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
					"spec":         "payg_version",
					"cfw_log":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"spec":         "payg_version",
						"cfw_log":      "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstance_basic1(t *testing.T) {
	ensureCloudFirewallInstanceTestRegion()
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "QueryAvailableInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccCloudFirewallInstancePreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"spec":         "enterprise_version",
					"ip_number":    "50",
					"band_width":   "50",
					"cfw_log":      "false",
					"logistics":    cloudFirewallInstanceLogisticsConfig,
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"spec":         "enterprise_version",
						"ip_number":    "50",
						"cfw_log":      "false",
						"logistics":    cloudFirewallInstanceLogisticsState,
						"period":       "1",
					}),
				),
			},
			// premium_version does not support fw_vpc_number
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"fw_vpc_number": "3",
			//		"modify_type":   "Upgrade",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"fw_vpc_number": "3",
			//			"modify_type":   "Upgrade",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":          "ultimate_version",
					"ip_number":     "400",
					"band_width":    "200",
					"fw_vpc_number": "5",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":          "ultimate_version",
						"ip_number":     "400",
						"fw_vpc_number": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":         "true",
					"cfw_log_storage": "5000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log":         "true",
						"cfw_log_storage": "5000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_number": "405",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_number": "405",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fw_vpc_number": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fw_vpc_number": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Month",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renew_period":          "1",
						"renewal_duration_unit": "Month",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "ManualRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status":        "ManualRenewal",
						"renewal_duration":      acctest.REMOVEKEY,
						"renew_period":          acctest.REMOVEKEY,
						"renewal_duration_unit": acctest.REMOVEKEY,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"renewal_duration": acctest.REMOVEKEY,
			//		"renew_period":     "2",
			//		"renewal_status":   "AutoRenewal",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"renewal_duration":      "2",
			//			"renew_period":          "2",
			//			"renewal_duration_unit": "Month",
			//			"renewal_status":        "AutoRenewal",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"renewal_status": "NotRenewal",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"renewal_status":        "NotRenewal",
			//			"renewal_duration":      acctest.REMOVEKEY,
			//			"renew_period":          acctest.REMOVEKEY,
			//			"renewal_duration_unit": acctest.REMOVEKEY,
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"band_width", "period", "modify_type", "logistics"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallInstance_unsupportedOrderParameters(t *testing.T) {
	ensureCloudFirewallInstanceTestRegion()
	resourceId := "alicloud_cloud_firewall_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccCloudFirewallInstancePreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"spec":           "enterprise_version",
					"ip_number":      "50",
					"band_width":     "50",
					"cfw_log":        "false",
					"cfw_account":    "true",
					"account_number": "1",
					"instance_count": "5",
					"logistics":      cloudFirewallInstanceLogisticsConfig,
					"period":         "1",
				}),
				ExpectError: regexp.MustCompile("InvalidParameter"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_account":    "false",
					"account_number": "2",
					"instance_count": "10",
				}),
				ExpectError: regexp.MustCompile("InvalidParameter"),
			},
		},
	})
}

var AliCloudCloudFirewallInstanceMap0 = map[string]string{
	"user_status": acctest.CHECKSET,
	"status":      acctest.CHECKSET,
}

func ensureCloudFirewallInstanceTestRegion() {
	if acctest.DefaultRegionToTest != "" {
		return
	}
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		acctest.DefaultRegionToTest = v
		return
	}
	acctest.DefaultRegionToTest = "cn-beijing"
	os.Setenv("ALICLOUD_REGION", acctest.DefaultRegionToTest)
}

func testAccCloudFirewallInstancePreCheck(t *testing.T) {
	acctest.PreCheck(t)
	ensureCloudFirewallInstanceTestRegion()
	testAccPreCheckForCleanUpInstances(t, "", "vipcloudfw", "vipcloudfw", "cfw", "cfw_pre_intl")
}

func AliCloudCloudFirewallInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}
`, name)
}
