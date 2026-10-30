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

func TestAccAliCloudGaBandwidthPackageAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackageAttachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${data.alicloud_ga_accelerators.default.ids.0}",
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"bandwidth_package_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": acctest.CHECKSET,
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

func TestAccAliCloudGaBandwidthPackageAttachment_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackageAttachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${data.alicloud_ga_accelerators.default.ids.0}",
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"bandwidth_package_id": acctest.CHECKSET,
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

func TestAccAliCloudGaBandwidthPackageAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackageAttachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageAttachmentBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${data.alicloud_ga_basic_accelerators.default.ids.0}",
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"bandwidth_package_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": acctest.CHECKSET,
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

func TestAccAliCloudGaBandwidthPackageAttachment_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackageAttachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageAttachmentBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${data.alicloud_ga_basic_accelerators.default.ids.0}",
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"bandwidth_package_id": acctest.CHECKSET,
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

var AliCloudGaBandwidthPackageAttachmentMap0 = map[string]string{
	"accelerators.#": acctest.CHECKSET,
	"status":         acctest.CHECKSET,
}

func AliCloudGaBandwidthPackageAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status                 = "active"
  		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth       = 100
  		type            = "Basic"
  		bandwidth_type  = "Basic"
  		payment_type    = "PayAsYouGo"
  		billing_type    = "PayBy95"
  		ratio           = 30
  		auto_pay        = true
  		auto_use_coupon = true
	}

	resource "alicloud_ga_bandwidth_package" "update" {
  		bandwidth       = 100
  		type            = "Basic"
  		bandwidth_type  = "Basic"
  		payment_type    = "PayAsYouGo"
  		billing_type    = "PayBy95"
  		ratio           = 30
  		auto_pay        = true
  		auto_use_coupon = true
	}
`, name)
}

func AliCloudGaBandwidthPackageAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_basic_accelerators" "default" {
  		status                 = "active"
  		bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth      = 20
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		duration       = 1
  		auto_pay       = true
	}

	resource "alicloud_ga_bandwidth_package" "update" {
  		bandwidth      = 20
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		duration       = 1
  		auto_pay       = true
	}
`, name)
}
