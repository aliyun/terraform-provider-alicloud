package ackone_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ackone"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass network.vswitches -- pre-existing ForceNew TypeList attribute; this change is a pure directory migration (alicloud/ -> alicloud/service/ackone/) that does not alter the vswitches schema, list handling or ordering semantics, so convergent reorder coverage is stock debt outside this change.

// Test AckOne Cluster. >>> Resource test cases, automatically generated.
// Case 4593
func TestAccAliCloudAckOneCluster_basic4593(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ackone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}",
							},
						},
					},
					"profile":        "Default",
					"argocd_enabled": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":   acctest.CHECKSET,
						"profile":        "Default",
						"argocd_enabled": "false",
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

func TestAccAliCloudAckOneCluster_basic4593_XFlow(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ackone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"profile": "XFlow",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": acctest.CHECKSET,
						"profile":      "XFlow",
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

func TestAccAliCloudAckOneCluster_basic4593_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ackone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"argocd_enabled": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

func TestAccAliCloudAckOneCluster_updateArgocdEnabled(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAckOneClusterMap4593)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ackone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sackonecluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAckOneClusterBasicDependence4593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.AckOneSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"profile":        "Default",
					"argocd_enabled": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":   acctest.CHECKSET,
						"profile":        "Default",
						"argocd_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name,
					"network": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
							"vswitches": []string{
								"${data.alicloud_vswitches.default.ids.0}"},
						},
					},
					"profile":        "Default",
					"argocd_enabled": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":   acctest.CHECKSET,
						"profile":        "Default",
						"argocd_enabled": "false",
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

var AlicloudAckOneClusterMap4593 = map[string]string{
	"status":       acctest.CHECKSET,
	"create_time":  acctest.CHECKSET,
	"cluster_name": acctest.CHECKSET,
	"profile":      acctest.CHECKSET,
}

// use existing vpc and vswitch because https://project.aone.alibaba-inc.com/v2/project/206563/bug/63062057
func AlicloudAckOneClusterBasicDependence4593(name string) string {
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

resource "alicloud_security_group" "default" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.default.ids.0
  security_group_name = "${var.name}"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
`, name)
}

// Test AckOne Cluster. <<< Resource test cases, automatically generated.
