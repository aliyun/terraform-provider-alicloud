// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order stock-debt declarations (pure rename migration; see scripts/collection-order/README.md)
// collection-order: bypass delivery_channel_filter.resource_types -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259

// Test ResourceManager DeliveryChannel. >>> Resource test cases, automatically generated.
// Case 当前账号投递渠道 11281
func TestAccAliCloudResourceManagerDeliveryChannel_basic11281(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_delivery_channel.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerDeliveryChannelMap11281)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerDeliveryChannel")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerDeliveryChannelBasicDependence11281)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls",
							"target_type": "SLS",
						},
					},
					"delivery_channel_name":        name,
					"delivery_channel_description": "delivery_channel_resource_spec_test",
					"delivery_channel_filter": []map[string]interface{}{
						{
							"resource_types": []string{
								"ACS::ECS::Instance", "ACS::ECS::Disk", "ACS::VPC::VPC"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "16:00Z",
							"target_arn":        "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls",
							"target_type":       "SLS",
							"custom_expression": "select * from resources limit 10;",
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":        name,
						"delivery_channel_description": "delivery_channel_resource_spec_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-1",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-1",
							"target_type": "SLS",
						},
					},
					"delivery_channel_name":        name + "_update",
					"delivery_channel_description": "delivery_channel_resource_spec_test_2",
					"delivery_channel_filter": []map[string]interface{}{
						{
							"resource_types": []string{
								"ACS::ECS::Instance", "ACS::VPC::VPC"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "17:00Z",
							"target_arn":        "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-1",
							"target_type":       "OSS",
							"custom_expression": "select * from resources;",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":        name + "_update",
						"delivery_channel_description": "delivery_channel_resource_spec_test_2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-2",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-2",
							"target_type": "SLS",
							"enabled":     "true",
						},
					},
					"delivery_channel_name":        name,
					"delivery_channel_description": "delivery_channel_resource_spec_test_3",
					"delivery_channel_filter": []map[string]interface{}{
						{
							"resource_types": []string{
								"ACS::ACK::Cluster", "ACS::Ons::Instance"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "21:00Z",
							"target_arn":        "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-2",
							"target_type":       "SLS",
							"custom_expression": "select * from resources limit 100;",
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-2",
								},
							},
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":        name,
						"delivery_channel_description": "delivery_channel_resource_spec_test_3",
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

var AlicloudResourceManagerDeliveryChannelMap11281 = map[string]string{}

func AlicloudResourceManagerDeliveryChannelBasicDependence11281(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 资源投递渠道测试 6316
func TestAccAliCloudResourceManagerDeliveryChannel_basic6316(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_delivery_channel.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerDeliveryChannelMap6316)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerDeliveryChannel")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerDeliveryChannelBasicDependence6316)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls",
							"target_type": "SLS",
						},
					},
					"delivery_channel_name":        name,
					"delivery_channel_description": "123",
					"delivery_channel_filter": []map[string]interface{}{
						{
							"resource_types": []string{
								"ALL"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "16:00Z",
							"target_arn":        "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls",
							"target_type":       "SLS",
							"custom_expression": "select * from resources limit 10;",
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":        name,
						"delivery_channel_description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-1",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-1",
							"target_type": "SLS",
						},
					},
					"delivery_channel_name":        name + "_update",
					"delivery_channel_description": "821",
					"delivery_channel_filter": []map[string]interface{}{
						{
							"resource_types": []string{
								"ACS::VPC::VPC"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "17:00Z",
							"target_arn":        "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-1",
							"target_type":       "OSS",
							"custom_expression": "select * from resources;",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":        name + "_update",
						"delivery_channel_description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-2",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-2",
							"target_type": "SLS",
							"enabled":     "true",
						},
					},
					"delivery_channel_name": name,
					"delivery_channel_filter": []map[string]interface{}{
						{
							"resource_types": []string{
								"ALL"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "21:00Z",
							"target_arn":        "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-2",
							"target_type":       "SLS",
							"custom_expression": "select * from resources limit 100;",
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-2",
								},
							},
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
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

var AlicloudResourceManagerDeliveryChannelMap6316 = map[string]string{}

func AlicloudResourceManagerDeliveryChannelBasicDependence6316(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ResourceManager DeliveryChannel. <<< Resource test cases, automatically generated.
