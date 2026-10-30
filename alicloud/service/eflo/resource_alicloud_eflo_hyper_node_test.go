// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass data_disk -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo HyperNode. >>> Resource test cases, automatically generated.
// Case 超节点资源测试 11893
func TestAccAliCloudEfloHyperNode_basic11893(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_hyper_node.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloHyperNodeMap11893)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloHyperNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloHyperNodeBasicDependence11893)
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
						"stage_num":        acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
					"payment_type": acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudEfloHyperNodeBasicDependence11893(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Eflo HyperNode. <<< Resource test cases, automatically generated.
