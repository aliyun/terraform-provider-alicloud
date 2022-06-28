package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCTrafficMirrorFilterIngressRulesDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_traffic_mirror_filter_ingress_rules.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	name := fmt.Sprintf("tf-testacc-vpctrafficmirrorfilteringressrule-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcTrafficMirrorFilterIngressRulesDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.id}",
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_filter_ingress_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.id}",
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_filter_ingress_rule.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.id}",
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_filter_ingress_rule.default.id}"},
			"status":                   "Created",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"traffic_mirror_filter_id": "${alicloud_vpc_traffic_mirror_filter.default.id}",
			"ids":                      []string{"${alicloud_vpc_traffic_mirror_filter_ingress_rule.default.id}"},
			"status":                   "Deleting",
		}),
	}
	var existActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"rules.#":                        "1",
			"rules.0.status":                 "Created",
			"rules.0.priority":               "1",
			"rules.0.rule_action":            "accept",
			"rules.0.protocol":               "UDP",
			"rules.0.destination_cidr_block": "10.0.0.0/24",
			"rules.0.source_cidr_block":      "10.0.0.0/24",
			"rules.0.destination_port_range": "1/120",
			"rules.0.source_port_range":      "1/120",
		}
	}

	var fakeActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"rules.#": "0",
		}
	}

	var vpcTrafficMirrorFilterIngressRuleCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existActiontrailTrailMapFunc,
		fakeMapFunc:  fakeActiontrailTrailMapFunc,
	}

	vpcTrafficMirrorFilterIngressRuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf)
}

func dataSourceVpcTrafficMirrorFilterIngressRulesDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}
resource "alicloud_vpc_traffic_mirror_filter" "default" {
  traffic_mirror_filter_name = var.name
}
resource "alicloud_vpc_traffic_mirror_filter_ingress_rule" "default" {
  traffic_mirror_filter_id = alicloud_vpc_traffic_mirror_filter.default.id
  priority                 = "1"
  rule_action              = "accept"
  protocol                 = "UDP"
  destination_cidr_block   = "10.0.0.0/24"
  source_cidr_block        = "10.0.0.0/24"
  destination_port_range   = "1/120"
  source_port_range        = "1/120"
}`, name)
}
