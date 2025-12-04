// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Mongodb Node. >>> Resource test cases, automatically generated.
// Case node子资源测试_mongos节点 6181
func TestAccAliCloudMongodbNode_basic6181(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_node.default"
	ra := resourceAttrInit(resourceId, AlicloudMongodbNodeMap6181)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongodbNodeBasicDependence6181)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_class":        "mdb.shard.4x.large.d",
					"readonly_replicas": "0",
					"db_instance_id":    "${alicloud_mongodb_instance.default4kU8ZI.id}",
					"node_type":         "mongos",
					"account_name":      "root",
					"account_password":  "q1w2e3r4!",
					"shard_direct":      "false",
					"business_info":     "1234",
					"auto_pay":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_class":        "mdb.shard.4x.large.d",
						"readonly_replicas": "0",
						"db_instance_id":    CHECKSET,
						"node_type":         "mongos",
						"account_name":      "root",
						"account_password":  "q1w2e3r4!",
						"shard_direct":      "false",
						"business_info":     CHECKSET,
						"auto_pay":          "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_class":     "mdb.shard.8x.large.d",
					"business_info":  "{“ActivityId\\\":\\\"000000000\\\"}",
					"order_type":     "UPGRADE",
					"effective_time": "Immediately",
					"from_app":       "OpenApi",
					"switch_time":    "2022-01-05T03:18:53Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_class":     "mdb.shard.8x.large.d",
						"business_info":  "{“ActivityId\":\"000000000\"}",
						"order_type":     "UPGRADE",
						"effective_time": "Immediately",
						"from_app":       "OpenApi",
						"switch_time":    "2022-01-05T03:18:53Z",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_name", "account_password", "auto_pay", "business_info", "effective_time", "from_app", "order_type", "shard_direct", "switch_time"},
			},
		},
	})
}

var AlicloudMongodbNodeMap6181 = map[string]string{
	"status":  CHECKSET,
	"node_id": CHECKSET,
}

func AlicloudMongodbNodeBasicDependence6181(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "ipv4网段-b" {
  default = "10.0.0.0/24"
}

resource "alicloud_vpc" "defaultFBPRtb" {
  description = "bgg-test"
  vpc_name    = "bgg-vpc-shanghai-b"
  cidr_block  = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaulti0SRXN" {
  vpc_id       = alicloud_vpc.defaultFBPRtb.id
  zone_id      = var.zone_id
  cidr_block   = var.ipv4网段-b
  vswitch_name = "bgg-shanghai-B"
}

resource "alicloud_mongodb_instance" "default4kU8ZI" {
  zone_id        = var.zone_id
  engine_version = "4.4"
  storage_type   = "cloud_essd1"
  vswitch_id     = alicloud_vswitch.defaulti0SRXN.id
  vpc_id         = alicloud_vpc.defaultFBPRtb.id
  db_instance_storage = "20"
  db_instance_class   = "mdb.shard.4x.large.d"
}


`, name)
}

// Case node子资源测试_shards节点 6201
func TestAccAliCloudMongodbNode_basic6201(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_node.default"
	ra := resourceAttrInit(resourceId, AlicloudMongodbNodeMap6201)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongodbNodeBasicDependence6201)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_class":        "mdb.shard.4x.large.d",
					"readonly_replicas": "0",
					"db_instance_id":    "${alicloud_mongodb_instance.default4kU8ZI.id}",
					"node_type":         "shard",
					"account_name":      "root",
					"account_password":  "q1w2e3r4!",
					"shard_direct":      "false",
					"business_info":     "1234",
					"auto_pay":          "true",
					"node_storage":      "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_class":        "mdb.shard.4x.large.d",
						"readonly_replicas": "0",
						"db_instance_id":    CHECKSET,
						"node_type":         "shard",
						"account_name":      "root",
						"account_password":  "q1w2e3r4!",
						"shard_direct":      "false",
						"business_info":     CHECKSET,
						"auto_pay":          "true",
						"node_storage":      "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_class":        "mdb.shard.8x.large.d",
					"readonly_replicas": "1",
					"node_storage":      "40",
					"order_type":        "UPGRADE",
					"effective_time":    "Immediately",
					"from_app":          "OpenApi",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_class":        "mdb.shard.8x.large.d",
						"readonly_replicas": "1",
						"node_storage":      "40",
						"order_type":        "UPGRADE",
						"effective_time":    "Immediately",
						"from_app":          "OpenApi",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_name", "account_password", "auto_pay", "business_info", "effective_time", "from_app", "order_type", "shard_direct", "switch_time"},
			},
		},
	})
}

var AlicloudMongodbNodeMap6201 = map[string]string{
	"status":  CHECKSET,
	"node_id": CHECKSET,
}

func AlicloudMongodbNodeBasicDependence6201(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "ipv4网段-b" {
  default = "10.0.0.0/24"
}

resource "alicloud_vpc" "defaultvAkfWC" {
  vpc_name   = "bgg-vpc-shanghai-b"
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultd9fTxe" {
  vpc_id       = alicloud_vpc.defaultvAkfWC.id
  zone_id      = var.zone_id
  cidr_block   = var.ipv4网段-b
  vswitch_name = "bgg-shanghai-B"
}

resource "alicloud_mongodb_instance" "default4kU8ZI" {
  zone_id        = var.zone_id
  engine_version = "4.4"
  storage_type   = "cloud_essd1"
  vswitch_id     = alicloud_vswitch.defaultd9fTxe.id
  vpc_id         = alicloud_vpc.defaultvAkfWC.id
}


`, name)
}

// Test Mongodb Node. <<< Resource test cases, automatically generated.
