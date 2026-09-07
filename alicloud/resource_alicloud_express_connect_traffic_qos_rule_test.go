package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AlicloudExpressConnectTrafficQosRuleMap6833 = map[string]string{
	"status":  CHECKSET,
	"rule_id": CHECKSET,
}

func AlicloudExpressConnectTrafficQosRuleBasicDependence6833(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_traffic_qos" "createQos" {
  qos_name        = "meijian-test"
  qos_description = "meijian-test"
}

resource "alicloud_express_connect_traffic_qos_association" "associateQos" {
  instance_id   = data.alicloud_express_connect_physical_connections.default.ids.1
  qos_id        = alicloud_express_connect_traffic_qos.createQos.id
  instance_type = "PHYSICALCONNECTION"
}

resource "alicloud_express_connect_traffic_qos_queue" "createQosQueue" {
  qos_id            = alicloud_express_connect_traffic_qos.createQos.id
  bandwidth_percent = "60"
  queue_description = "meijian-test"
  queue_name        = "meijian-test"
  queue_type        = "Medium"
}


`, name)
}

// Case Qos规则-v6-线上 6834
func TestAccAliCloudExpressConnectTrafficQosRule_basic6834(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosRuleMap6834)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosRuleBasicDependence6834)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_id": "${alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id}",
					"qos_id":   "${alicloud_express_connect_traffic_qos.createQos.id}",
					"priority": "1",
					"protocol": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_id": CHECKSET,
						"qos_id":   CHECKSET,
						"priority": "1",
						"protocol": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp": "-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp": "-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_port_range": "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_port_range": "-1/-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remarking_dscp": "-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remarking_dscp": "-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"src_port_range": "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_port_range": "-1/-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_ipv6_cidr": "2001:db8:1234:5678::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_ipv6_cidr": "2001:db8:1234:5678::/64",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"src_ipv6_cidr": "2001:db8:1234:5678::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_ipv6_cidr": "2001:db8:1234:5678::/64",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "meijian-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "meijian-test-1",
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
					"remarking_dscp": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remarking_dscp": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "ICMP(IPv4)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "ICMP(IPv4)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_ipv6_cidr": "2001:db8:1234:5679::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_ipv6_cidr": "2001:db8:1234:5679::/64",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"src_ipv6_cidr": "2001:db8:1234:5679::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_ipv6_cidr": "2001:db8:1234:5679::/64",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test",
					"queue_id":         "${alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id}",
					"qos_id":           "${alicloud_express_connect_traffic_qos.createQos.id}",
					"rule_name":        "meijian-test",
					"match_dscp":       "-1",
					"priority":         "1",
					"dst_port_range":   "-1/-1",
					"remarking_dscp":   "-1",
					"protocol":         "ALL",
					"src_port_range":   "-1/-1",
					"dst_ipv6_cidr":    "2001:db8:1234:5678::/64",
					"src_ipv6_cidr":    "2001:db8:1234:5678::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test",
						"queue_id":         CHECKSET,
						"qos_id":           CHECKSET,
						"rule_name":        "meijian-test",
						"match_dscp":       "-1",
						"priority":         "1",
						"dst_port_range":   "-1/-1",
						"remarking_dscp":   "-1",
						"protocol":         "ALL",
						"src_port_range":   "-1/-1",
						"dst_ipv6_cidr":    "2001:db8:1234:5678::/64",
						"src_ipv6_cidr":    "2001:db8:1234:5678::/64",
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

var AlicloudExpressConnectTrafficQosRuleMap6834 = map[string]string{
	"status":  CHECKSET,
	"rule_id": CHECKSET,
}

func AlicloudExpressConnectTrafficQosRuleBasicDependence6834(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_traffic_qos" "createQos" {
  qos_name        = "meijian-test"
  qos_description = "meijian-test"
}

resource "alicloud_express_connect_traffic_qos_association" "associateQos" {
  instance_id   = data.alicloud_express_connect_physical_connections.default.ids.1
  qos_id        = alicloud_express_connect_traffic_qos.createQos.id
  instance_type = "PHYSICALCONNECTION"
}

resource "alicloud_express_connect_traffic_qos_queue" "createQosQueue" {
  qos_id            = alicloud_express_connect_traffic_qos.createQos.id
  bandwidth_percent = "60"
  queue_description = "meijian-test"
  queue_name        = "meijian-test"
  queue_type        = "Medium"
}


`, name)
}

// Case Qos规则-v4-线上 6833  twin
func TestAccAliCloudExpressConnectTrafficQosRule_basic6833_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosRuleMap6833)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosRuleBasicDependence6833)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test",
					"queue_id":         "${alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id}",
					"qos_id":           "${alicloud_express_connect_traffic_qos.createQos.id}",
					"rule_name":        "meijian-test",
					"match_dscp":       "-1",
					"priority":         "1",
					"dst_port_range":   "-1/-1",
					"remarking_dscp":   "-1",
					"src_cidr":         "192.168.1.0/24",
					"protocol":         "ALL",
					"src_port_range":   "-1/-1",
					"dst_cidr":         "192.168.2.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test",
						"queue_id":         CHECKSET,
						"qos_id":           CHECKSET,
						"rule_name":        "meijian-test",
						"match_dscp":       "-1",
						"priority":         "1",
						"dst_port_range":   "-1/-1",
						"remarking_dscp":   "-1",
						"src_cidr":         "192.168.1.0/24",
						"protocol":         "ALL",
						"src_port_range":   "-1/-1",
						"dst_cidr":         "192.168.2.0/24",
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

// Case Qos规则-v6-线上 6834  twin
func TestAccAliCloudExpressConnectTrafficQosRule_basic6834_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosRuleMap6834)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosRuleBasicDependence6834)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test",
					"queue_id":         "${alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id}",
					"qos_id":           "${alicloud_express_connect_traffic_qos.createQos.id}",
					"rule_name":        "meijian-test",
					"match_dscp":       "-1",
					"priority":         "1",
					"dst_port_range":   "-1/-1",
					"remarking_dscp":   "-1",
					"protocol":         "ALL",
					"src_port_range":   "-1/-1",
					"dst_ipv6_cidr":    "2001:db8:1234:5678::/64",
					"src_ipv6_cidr":    "2001:db8:1234:5678::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test",
						"queue_id":         CHECKSET,
						"qos_id":           CHECKSET,
						"rule_name":        "meijian-test",
						"match_dscp":       "-1",
						"priority":         "1",
						"dst_port_range":   "-1/-1",
						"remarking_dscp":   "-1",
						"protocol":         "ALL",
						"src_port_range":   "-1/-1",
						"dst_ipv6_cidr":    "2001:db8:1234:5678::/64",
						"src_ipv6_cidr":    "2001:db8:1234:5678::/64",
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

// Case Qos规则-v4-线上 6833  raw
func TestAccAliCloudExpressConnectTrafficQosRule_basic6833_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosRuleMap6833)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosRuleBasicDependence6833)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test",
					"queue_id":         "${alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id}",
					"qos_id":           "${alicloud_express_connect_traffic_qos.createQos.id}",
					"rule_name":        "meijian-test",
					"match_dscp":       "-1",
					"priority":         "1",
					"dst_port_range":   "-1/-1",
					"remarking_dscp":   "-1",
					"src_cidr":         "192.168.1.0/24",
					"protocol":         "ALL",
					"src_port_range":   "-1/-1",
					"dst_cidr":         "192.168.2.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test",
						"queue_id":         CHECKSET,
						"qos_id":           CHECKSET,
						"rule_name":        "meijian-test",
						"match_dscp":       "-1",
						"priority":         "1",
						"dst_port_range":   "-1/-1",
						"remarking_dscp":   "-1",
						"src_cidr":         "192.168.1.0/24",
						"protocol":         "ALL",
						"src_port_range":   "-1/-1",
						"dst_cidr":         "192.168.2.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test-1",
					"rule_name":        "meijian-test-1",
					"priority":         "2",
					"remarking_dscp":   "1",
					"src_cidr":         "192.168.3.0/24",
					"protocol":         "ICMP(IPv4)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test-1",
						"rule_name":        "meijian-test-1",
						"priority":         "2",
						"remarking_dscp":   "1",
						"src_cidr":         "192.168.3.0/24",
						"protocol":         "ICMP(IPv4)",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_port_range": "1/1",
					"remarking_dscp": "2",
					"protocol":       "TCP",
					"src_port_range": "2/2",
					"dst_cidr":       "192.168.4.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_port_range": "1/1",
						"remarking_dscp": "2",
						"protocol":       "TCP",
						"src_port_range": "2/2",
						"dst_cidr":       "192.168.4.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_port_range": "3/3",
					"remarking_dscp": "-1",
					"protocol":       "UDP",
					"src_port_range": "4/4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_port_range": "3/3",
						"remarking_dscp": "-1",
						"protocol":       "UDP",
						"src_port_range": "4/4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp":     "1",
					"dst_port_range": "-1/-1",
					"protocol":       "GRE",
					"src_port_range": "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp":     "1",
						"dst_port_range": "-1/-1",
						"protocol":       "GRE",
						"src_port_range": "-1/-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp":     "2",
					"protocol":       "SSH",
					"dst_port_range": "22/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp":     "2",
						"protocol":       "SSH",
						"dst_port_range": "22/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp":     "-1",
					"protocol":       "Telnet",
					"dst_port_range": "23/23",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp":     "-1",
						"protocol":       "Telnet",
						"dst_port_range": "23/23",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp":     "1",
					"remarking_dscp": "1",
					"protocol":       "HTTP",
					"dst_port_range": "80/80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp":     "1",
						"remarking_dscp": "1",
						"protocol":       "HTTP",
						"dst_port_range": "80/80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp":     "2",
					"remarking_dscp": "2",
					"protocol":       "HTTPS",
					"dst_port_range": "443/443",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp":     "2",
						"remarking_dscp": "2",
						"protocol":       "HTTPS",
						"dst_port_range": "443/443",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_dscp":     "-1",
					"remarking_dscp": "-1",
					"protocol":       "MS SQL",
					"dst_port_range": "1443/1443",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_dscp":     "-1",
						"remarking_dscp": "-1",
						"protocol":       "MS SQL",
						"dst_port_range": "1443/1443",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":       "Oracle",
					"dst_port_range": "1521/1521",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":       "Oracle",
						"dst_port_range": "1521/1521",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":       "MySql",
					"dst_port_range": "3306/3306",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":       "MySql",
						"dst_port_range": "3306/3306",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":       "RDP",
					"dst_port_range": "3389/3389",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":       "RDP",
						"dst_port_range": "3389/3389",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":       "Postgre SQL",
					"dst_port_range": "5432/5432",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":       "Postgre SQL",
						"dst_port_range": "5432/5432",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":       "Redis",
					"dst_port_range": "6379/6379",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":       "Redis",
						"dst_port_range": "6379/6379",
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

// Case Qos规则-v6-线上 6834  raw
func TestAccAliCloudExpressConnectTrafficQosRule_basic6834_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosRuleMap6834)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosRuleBasicDependence6834)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test",
					"queue_id":         "${alicloud_express_connect_traffic_qos_queue.createQosQueue.queue_id}",
					"qos_id":           "${alicloud_express_connect_traffic_qos.createQos.id}",
					"rule_name":        "meijian-test",
					"match_dscp":       "-1",
					"priority":         "1",
					"dst_port_range":   "-1/-1",
					"remarking_dscp":   "-1",
					"protocol":         "ALL",
					"src_port_range":   "-1/-1",
					"dst_ipv6_cidr":    "2001:db8:1234:5678::/64",
					"src_ipv6_cidr":    "2001:db8:1234:5678::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test",
						"queue_id":         CHECKSET,
						"qos_id":           CHECKSET,
						"rule_name":        "meijian-test",
						"match_dscp":       "-1",
						"priority":         "1",
						"dst_port_range":   "-1/-1",
						"remarking_dscp":   "-1",
						"protocol":         "ALL",
						"src_port_range":   "-1/-1",
						"dst_ipv6_cidr":    "2001:db8:1234:5678::/64",
						"src_ipv6_cidr":    "2001:db8:1234:5678::/64",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_description": "meijian-test-1",
					"rule_name":        "meijian-test-1",
					"priority":         "2",
					"remarking_dscp":   "1",
					"protocol":         "ICMP(IPv4)",
					"dst_ipv6_cidr":    "2001:db8:1234:5679::/64",
					"src_ipv6_cidr":    "2001:db8:1234:5679::/64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_description": "meijian-test-1",
						"rule_name":        "meijian-test-1",
						"priority":         "2",
						"remarking_dscp":   "1",
						"protocol":         "ICMP(IPv4)",
						"dst_ipv6_cidr":    "2001:db8:1234:5679::/64",
						"src_ipv6_cidr":    "2001:db8:1234:5679::/64",
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
