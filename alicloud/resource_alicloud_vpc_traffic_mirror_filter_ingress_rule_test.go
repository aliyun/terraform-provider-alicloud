package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVPCTrafficMirrorFilterIngressRule_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	resourceId := "alicloud_vpc_traffic_mirror_filter_ingress_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCTrafficMirrorFilterIngressRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcTrafficMirrorFilterIngressRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorfilteringressrule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCTrafficMirrorFilterIngressRuleBasicDependence0)
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
					"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.id}",
					"priority":                 "1",
					"rule_action":              "accept",
					"protocol":                 "ICMP",
					"destination_cidr_block":   "10.0.0.0/24",
					"source_cidr_block":        "10.0.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_id": CHECKSET,
						"priority":                 "1",
						"rule_action":              "accept",
						"protocol":                 "ICMP",
						"destination_cidr_block":   "10.0.0.0/24",
						"source_cidr_block":        "10.0.0.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "UDP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "UDP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_cidr_block": "10.0.0.0/20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidr_block": "10.0.0.0/20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_port_range": "1/120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_port_range": "1/120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_cidr_block": "10.0.0.0/20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_block": "10.0.0.0/20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_port_range": "1/120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_port_range": "1/120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_action": "drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_action": "drop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":               "1",
					"rule_action":            "accept",
					"protocol":               "ICMP",
					"destination_cidr_block": "10.0.0.0/24",
					"source_cidr_block":      "10.0.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":               "1",
						"rule_action":            "accept",
						"protocol":               "ICMP",
						"destination_cidr_block": "10.0.0.0/24",
						"source_cidr_block":      "10.0.0.0/24",
						"source_port_range":      "-1/-1",
						"destination_port_range": "-1/-1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudVPCTrafficMirrorFilterIngressRule_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	resourceId := "alicloud_vpc_traffic_mirror_filter_ingress_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCTrafficMirrorFilterIngressRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcTrafficMirrorFilterIngressRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorfilteringressrule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCTrafficMirrorFilterIngressRuleBasicDependence0)
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
					"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.id}",
					"priority":                 "1",
					"rule_action":              "accept",
					"protocol":                 "UDP",
					"destination_cidr_block":   "10.0.0.0/24",
					"source_cidr_block":        "10.0.0.0/24",
					"dry_run":                  "false",
					"destination_port_range":   "1/120",
					"source_port_range":        "1/120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_mirror_filter_id": CHECKSET,
						"priority":                 "1",
						"rule_action":              "accept",
						"protocol":                 "UDP",
						"destination_cidr_block":   "10.0.0.0/24",
						"source_cidr_block":        "10.0.0.0/24",
						"dry_run":                  "false",
						"destination_port_range":   "1/120",
						"source_port_range":        "1/120",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudVPCTrafficMirrorFilterIngressRuleMap0 = map[string]string{
	"dry_run": NOSET,
	"status":  CHECKSET,
}

func AlicloudVPCTrafficMirrorFilterIngressRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_vpc_traffic_mirror_filter" "default" {
  traffic_mirror_filter_name = var.name
}

`, name)
}
