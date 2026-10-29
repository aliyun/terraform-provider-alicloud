package ga_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaBasicAccelerator_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBasicAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBasicAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaBasicAccelerator-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBasicAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_pay": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_pay": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"basic_accelerator_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"basic_accelerator_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BasicAccelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "BasicAccelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle", "auto_use_coupon", "auto_renew", "auto_renew_duration", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaBasicAccelerator_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBasicAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBasicAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaBasicAccelerator-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBasicAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "BandwidthPackage",
					"payment_type":           "Subscription",
					"auto_pay":               "true",
					"duration":               "1",
					"pricing_cycle":          "Month",
					"auto_use_coupon":        "true",
					"auto_renew":             "false",
					"auto_renew_duration":    "1",
					"promotion_option_no":    "",
					"basic_accelerator_name": name,
					"description":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BasicAccelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "BandwidthPackage",
						"payment_type":           "Subscription",
						"basic_accelerator_name": name,
						"description":            name,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "BasicAccelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle", "auto_use_coupon", "auto_renew", "auto_renew_duration", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaBasicAccelerator_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBasicAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBasicAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaBasicAccelerator-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBasicAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"payment_type":           "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "CDT",
						"payment_type":           "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"basic_accelerator_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"basic_accelerator_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BasicAccelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "BasicAccelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle", "auto_use_coupon", "auto_renew", "auto_renew_duration", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaBasicAccelerator_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBasicAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBasicAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaBasicAccelerator-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBasicAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"payment_type":           "PayAsYouGo",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"basic_accelerator_name": name,
					"description":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BasicAccelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "CDT",
						"payment_type":           "PayAsYouGo",
						"resource_group_id":      acctest.CHECKSET,
						"basic_accelerator_name": name,
						"description":            name,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "BasicAccelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "pricing_cycle", "auto_use_coupon", "auto_renew", "auto_renew_duration", "promotion_option_no"},
			},
		},
	})
}

var AliCloudGaBasicAcceleratorMap0 = map[string]string{
	"payment_type":      acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
}

func AliCloudGaBasicAcceleratorBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`)
}
