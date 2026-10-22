// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo HyperNode. >>> Resource test cases, automatically generated.
// Case 超节点资源测试 11893
func TestAccAliCloudEfloHyperNode_basic11893(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_hyper_node.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloHyperNodeMap11893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloHyperNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloHyperNodeBasicDependence11893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":          "ap-southeast-7a",
					"machine_type":     "efg3.GN9A.ch72",
					"hpn_zone":         "A1",
					"server_arch":      "bmserver",
					"payment_duration": "1",
					"payment_type":     "Subscription",
					"stage_num":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":          "ap-southeast-7a",
						"machine_type":     "efg3.GN9A.ch72",
						"hpn_zone":         "A1",
						"server_arch":      "bmserver",
						"payment_duration": "1",
						"payment_type":     "Subscription",
						"stage_num":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hostname":          "tfacc-host",
					"cluster_id":        "ehc-tfacc-placeholder",
					"node_group_id":     "ng-tfacc-placeholder",
					"login_password":    "YourPassword_123",
					"renewal_duration":  "1",
					"renewal_status":    "AutoRenewal",
					"resource_group_id": "rg-tfacc-placeholder",
					"user_data":         "I3RmYWNjCg==",
					"vswitch_id":        "vsw-tfacc-placeholder",
					"vpc_id":            "vpc-tfacc-placeholder",
					"data_disk": []map[string]interface{}{
						{
							"bursting_enabled":  "true",
							"delete_with_node":  "true",
							"category":          "cloud_essd",
							"size":              "2000",
							"performance_level": "PL1",
							"provisioned_iops":  "50000",
						},
					},
					"payment_type": REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_disk", "login_password", "payment_duration", "stage_num"},
			},
		},
	})
}

var AlicloudEfloHyperNodeMap11893 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEfloHyperNodeBasicDependence11893(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Eflo HyperNode. <<< Resource test cases, automatically generated.
