package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudStorageGatewayGateway_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudStorageGatewayGatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCloudStorageGatewayGateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudStorageGatewayGatewayBasicDependence0)
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
					"storage_bundle_id": "${data.alicloud_cloud_storage_gateway_storage_bundles.default.bundles.0.id}",
					"type":              "Iscsi",
					"location":          "Cloud",
					"gateway_name":      name,
					"gateway_class":     "Basic",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_bundle_id": CHECKSET,
						"type":              "Iscsi",
						"location":          "Cloud",
						"gateway_name":      name,
						"gateway_class":     "Basic",
						"vswitch_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_class": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_class": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_class": "Enhanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_class": "Enhanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_class": "Advanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_class": "Advanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_network_bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_network_bandwidth": "10",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "reason_detail"},
			},
		},
	})
}

func TestAccAliCloudCloudStorageGatewayGateway_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudStorageGatewayGatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCloudStorageGatewayGateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudStorageGatewayGatewayBasicDependence0)
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
					"storage_bundle_id":        "${data.alicloud_cloud_storage_gateway_storage_bundles.default.bundles.0.id}",
					"type":                     "Iscsi",
					"location":                 "Cloud",
					"gateway_name":             name,
					"gateway_class":            "Basic",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"public_network_bandwidth": "10",
					"payment_type":             "PayAsYouGo",
					"description":              name,
					"reason_type":              "REASON2",
					"reason_detail":            name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_bundle_id":        CHECKSET,
						"type":                     "Iscsi",
						"location":                 "Cloud",
						"gateway_name":             name,
						"gateway_class":            "Basic",
						"vswitch_id":               CHECKSET,
						"public_network_bandwidth": "10",
						"payment_type":             "PayAsYouGo",
						"description":              name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "reason_detail"},
			},
		},
	})
}

func TestAccAliCloudCloudStorageGatewayGateway_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudStorageGatewayGatewayMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCloudStorageGatewayGateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudStorageGatewayGatewayBasicDependence0)
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
					"storage_bundle_id": "${data.alicloud_cloud_storage_gateway_storage_bundles.default.bundles.0.id}",
					"type":              "File",
					"location":          "On_Premise",
					"gateway_name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_bundle_id": CHECKSET,
						"type":              "File",
						"location":          "On_Premise",
						"gateway_name":      name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name + "-update",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "reason_detail"},
			},
		},
	})
}

func TestAccAliCloudCloudStorageGatewayGateway_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudStorageGatewayGatewayMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCloudStorageGatewayGateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudStorageGatewayGatewayBasicDependence0)
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
					"storage_bundle_id": "${data.alicloud_cloud_storage_gateway_storage_bundles.default.bundles.0.id}",
					"type":              "File",
					"location":          "On_Premise",
					"gateway_name":      name,
					"payment_type":      "PayAsYouGo",
					"description":       name,
					"reason_type":       "REASON2",
					"reason_detail":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_bundle_id": CHECKSET,
						"type":              "File",
						"location":          "On_Premise",
						"gateway_name":      name,
						"payment_type":      "PayAsYouGo",
						"description":       name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "reason_detail"},
			},
		},
	})
}

func TestAccAliCloudCloudStorageGatewayGateway_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudStorageGatewayGatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCloudStorageGatewayGateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudStorageGatewayGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_bundle_id": "${data.alicloud_cloud_storage_gateway_storage_bundles.default.bundles.0.id}",
					"type":              "Iscsi",
					"location":          "Cloud",
					"gateway_name":      name,
					"gateway_class":     "Basic",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"payment_type":      "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_bundle_id": CHECKSET,
						"type":              "Iscsi",
						"location":          "Cloud",
						"gateway_name":      name,
						"gateway_class":     "Basic",
						"vswitch_id":        CHECKSET,
						"payment_type":      "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name": name + "-update",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "reason_detail"},
			},
		},
	})
}

func TestAccAliCloudCloudStorageGatewayGateway_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudStorageGatewayGatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCloudStorageGatewayGateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudStorageGatewayGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_bundle_id":        "${data.alicloud_cloud_storage_gateway_storage_bundles.default.bundles.0.id}",
					"type":                     "Iscsi",
					"location":                 "Cloud",
					"gateway_name":             name,
					"gateway_class":            "Basic",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"public_network_bandwidth": "10",
					"payment_type":             "Subscription",
					"description":              name,
					"release_after_expiration": "true",
					"reason_type":              "REASON2",
					"reason_detail":            name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_bundle_id":        CHECKSET,
						"type":                     "Iscsi",
						"location":                 "Cloud",
						"gateway_name":             name,
						"gateway_class":            "Basic",
						"vswitch_id":               CHECKSET,
						"public_network_bandwidth": "10",
						"payment_type":             "Subscription",
						"description":              name,
						"release_after_expiration": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reason_type", "reason_detail"},
			},
		},
	})
}

var AliCloudCloudStorageGatewayGatewayMap0 = map[string]string{
	"public_network_bandwidth": CHECKSET,
	"status":                   CHECKSET,
}

var AliCloudCloudStorageGatewayGatewayMap1 = map[string]string{
	"status": CHECKSET,
}

func AliCloudCloudStorageGatewayGatewayBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
	}

	data "alicloud_cloud_storage_gateway_storage_bundles" "default" {
  		backend_bucket_region_id = "%s"
	}
`, name, defaultRegionToTest)
}
