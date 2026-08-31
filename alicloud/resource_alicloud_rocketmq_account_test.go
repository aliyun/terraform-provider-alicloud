package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rocketmq Account. >>> Resource test cases, automatically generated.
// Case account测试 10054
func TestAccAliCloudRocketmqAccount_basic10054(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_account.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqAccountMap10054)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrocketmq%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqAccountBasicDependence10054)
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
					"instance_id": "${alicloud_rocketmq_instance.default9hAb83.id}",
					"username":    name,
					"password":    "1739867871",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"username":    name,
						"password":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "1739867872",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_status": "DISABLE",
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

var AliCloudRocketmqAccountMap10054 = map[string]string{
	"account_status": CHECKSET,
}

func AliCloudRocketmqAccountBasicDependence10054(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "defaultg6ZXs2" {
  description = "111"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "pop-test-vpc"
}

resource "alicloud_vswitch" "defaultvMQbCy" {
  vpc_id       = alicloud_vpc.defaultg6ZXs2.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "pop-test-vswitch"
}

resource "alicloud_rocketmq_instance" "default9hAb83" {
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
  software {
    maintain_time = "02:00-06:00"
  }

  tags = {
    Created = "TF"
    For     = "example"
  }
  network_info {
    vpc_info {
      vpc_id     = alicloud_vpc.defaultg6ZXs2.id
      vswitch_id = alicloud_vswitch.defaultvMQbCy.id
    }

    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
      ip_whitelist = [
        "192.168.0.0/16",
        "10.10.0.0/16",
        "172.168.0.0/16"
      ]
    }
  }
  acl_info {
    default_vpc_auth_free = false
    acl_types             = ["default", "apache_acl"]
  }
}
`, name)
}

// Case account测试 10054 twin
func TestAccAliCloudRocketmqAccount_basic10054_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_account.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqAccountMap10054)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrocketmq%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqAccountBasicDependence10054)
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
					"account_status": "ENABLE",
					"instance_id":    "${alicloud_rocketmq_instance.default9hAb83.id}",
					"username":       name,
					"password":       "1739867871",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_status": "ENABLE",
						"instance_id":    CHECKSET,
						"username":       name,
						"password":       CHECKSET,
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

// Test Rocketmq Account. <<< Resource test cases, automatically generated.
