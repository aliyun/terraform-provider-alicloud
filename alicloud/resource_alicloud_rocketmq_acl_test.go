package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rocketmq Acl. >>> Resource test cases, automatically generated.
// Case acl测试 10135
func TestAccAliCloudRocketmqAcl_basic10135(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_acl.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqAclMap10135)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrocketmq%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqAclBasicDependence10135)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"actions":       []string{"Pub"},
					"instance_id":   "${alicloud_rocketmq_instance.defaultKJZNVM.id}",
					"username":      "${alicloud_rocketmq_account.defaultMeNlxe.username}",
					"resource_name": "${alicloud_rocketmq_topic.defaultVA0zog.topic_name}",
					"resource_type": "Topic",
					"decision":      "Deny",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"actions.#":     "1",
						"instance_id":   CHECKSET,
						"username":      CHECKSET,
						"resource_name": CHECKSET,
						"resource_type": "Topic",
						"decision":      "Deny",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"actions": []string{"Pub", "Sub"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"actions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"decision": "Allow",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"decision": "Allow",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelists": []string{"192.168.1.1", "192.168.2.2", "192.168.3.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelists.#": "3",
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

func TestAccAliCloudRocketmqAcl_basic10135_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_acl.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqAclMap10135)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrocketmq%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqAclBasicDependence10135)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"actions":       []string{"Pub"},
					"instance_id":   "${alicloud_rocketmq_instance.defaultKJZNVM.id}",
					"username":      "${alicloud_rocketmq_account.defaultMeNlxe.username}",
					"resource_name": "${alicloud_rocketmq_topic.defaultVA0zog.topic_name}",
					"resource_type": "Topic",
					"decision":      "Deny",
					"ip_whitelists": []string{"192.168.1.1", "192.168.2.2", "192.168.3.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"actions.#":       "1",
						"instance_id":     CHECKSET,
						"username":        CHECKSET,
						"resource_name":   CHECKSET,
						"resource_type":   "Topic",
						"decision":        "Deny",
						"ip_whitelists.#": "3",
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

var AliCloudRocketmqAclMap10135 = map[string]string{}

func AliCloudRocketmqAclBasicDependence10135(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "defaultrqDtGm" {
  description = "1111"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "pop-test-vpc"
}

resource "alicloud_vswitch" "defaultjUrTYm" {
  vpc_id       = alicloud_vpc.defaultrqDtGm.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "pop-test-vswitch"
}

resource "alicloud_rocketmq_instance" "defaultKJZNVM" {
  product_info {
    msg_process_spec       = "rmq.p2.4xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code    = "rmq"
  series_code     = "professional"
  payment_type    = "PayAsYouGo"
  instance_name   = var.name
  sub_series_code = "cluster_ha"
  remark          = "example"
  network_info {
    vpc_info {
      vpc_id     = alicloud_vpc.defaultrqDtGm.id
      vswitch_id = alicloud_vswitch.defaultjUrTYm.id
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "5"
    }
  }
  acl_info {
    default_vpc_auth_free = false
    acl_types             = ["default", "apache_acl"]
  }
}

resource "alicloud_rocketmq_account" "defaultMeNlxe" {
  account_status = "ENABLE"
  instance_id    = alicloud_rocketmq_instance.defaultKJZNVM.id
  username       = "zhenyuantest"
  password       = "123456"
}

resource "alicloud_rocketmq_topic" "defaultVA0zog" {
  instance_id  = alicloud_rocketmq_instance.defaultKJZNVM.id
  message_type = "NORMAL"
  topic_name   = "zhenyuantest"
}
`, name)
}

// Test Rocketmq Acl. <<< Resource test cases, automatically generated.
