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
// collection-order: bypass delivery_channel_filter.account_scopes -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259
// collection-order: bypass delivery_channel_filter.resource_types -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259

// Test ResourceManager MultiAccountDeliveryChannel. >>> Resource test cases, automatically generated.
// Case 多账号投递渠道 7679
func TestAccAliCloudResourceManagerMultiAccountDeliveryChannel_basic7679(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_multi_account_delivery_channel.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerMultiAccountDeliveryChannelMap7679)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerMultiAccountDeliveryChannel")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerMultiAccountDeliveryChannelBasicDependence7679)
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
					"delivery_channel_description":        "multi_delivery_channel_resource_spec_mq_test",
					"multi_account_delivery_channel_name": name,
					"delivery_channel_filter": []map[string]interface{}{
						{
							"account_scopes": []string{
								"${alicloud_resource_manager_folder.defaultuHQ8Cu.id}", "${alicloud_resource_manager_folder.defaultioI16p.id}", "${alicloud_resource_manager_folder.default55Uum4.id}"},
							"resource_types": []string{
								"ACS::ACK::Cluster", "ACS::ActionTrail::Trail", "ACS::BPStudio::Application"},
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
						"delivery_channel_description":        "multi_delivery_channel_resource_spec_mq_test",
						"multi_account_delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_change_delivery": []map[string]interface{}{
						{
							"enabled": "true",
							"sls_properties": []map[string]interface{}{
								{
									"oversized_data_oss_target_arn": "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-1",
								},
							},
							"target_arn":  "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-test/logstore/resourcecenter-delivery-aone-test-sls-1",
							"target_type": "SLS",
						},
					},
					"delivery_channel_description":        "multi_delivery_channel_resource_spec_mq_test_2",
					"multi_account_delivery_channel_name": name + "_update",
					"delivery_channel_filter": []map[string]interface{}{
						{
							"account_scopes": []string{
								"${alicloud_resource_manager_folder.defaultiEjEbe.id}", "${alicloud_resource_manager_folder.defaultdNL2TN.id}"},
							"resource_types": []string{
								"ACS::CBWP::CommonBandwidthPackage", "ACS::CDN::Domain"},
						},
					},
					"resource_snapshot_delivery": []map[string]interface{}{
						{
							"delivery_time":     "16:01Z",
							"enabled":           "true",
							"target_arn":        "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-test-delivery-oss-1",
							"target_type":       "OSS",
							"custom_expression": "select * from resources;",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_description":        "multi_delivery_channel_resource_spec_mq_test_2",
						"multi_account_delivery_channel_name": name + "_update",
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
					"delivery_channel_description":        "multi_delivery_channel_resource_spec_test_3",
					"multi_account_delivery_channel_name": name,
					"delivery_channel_filter": []map[string]interface{}{
						{
							"account_scopes": []string{
								"${alicloud_resource_manager_folder.default55Uum4.id}"},
							"resource_types": []string{
								"ACS::ECS::Instance", "ACS::VPC::VPC"},
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
						"delivery_channel_description":        "multi_delivery_channel_resource_spec_test_3",
						"multi_account_delivery_channel_name": name,
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

var AlicloudResourceManagerMultiAccountDeliveryChannelMap7679 = map[string]string{}

func AlicloudResourceManagerMultiAccountDeliveryChannelBasicDependence7679(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_resource_manager_folder" "defaultuHQ8Cu" {
  folder_name = "folder-aone-test-1"
}

resource "alicloud_resource_manager_folder" "defaultioI16p" {
  folder_name = "folder-aone-test-2"
}

resource "alicloud_resource_manager_folder" "default55Uum4" {
  folder_name = "folder-aone-test-3"
}

resource "alicloud_resource_manager_folder" "defaultiEjEbe" {
  folder_name = "folder-aone-test-4"
}

resource "alicloud_resource_manager_folder" "defaultdNL2TN" {
  folder_name = "folder-aone-test-5"
}


`, name)
}

// Test ResourceManager MultiAccountDeliveryChannel. <<< Resource test cases, automatically generated.
