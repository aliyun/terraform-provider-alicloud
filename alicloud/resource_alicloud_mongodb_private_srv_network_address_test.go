package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Mongodb PrivateSrvNetworkAddress. >>> Resource test cases, automatically generated.
// Case 私有网络srv测试 9657
func TestAccAliCloudMongodbPrivateSrvNetworkAddress_basic9657(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_private_srv_network_address.default"
	ra := resourceAttrInit(resourceId, AlicloudMongodbPrivateSrvNetworkAddressMap9657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbPrivateSrvNetworkAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smongodbprivatesrvnetworkaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongodbPrivateSrvNetworkAddressBasicDependence9657)
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
					"db_instance_id": "${alicloud_mongodb_instance.defaultHrZmxC.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
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

var AlicloudMongodbPrivateSrvNetworkAddressMap9657 = map[string]string{
	"private_srv_connection_string_uri": CHECKSET,
}

func AlicloudMongodbPrivateSrvNetworkAddressBasicDependence9657(name string) string {
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

resource "alicloud_vpc" "defaultie35CW" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultg0DCAR" {
  vpc_id     = alicloud_vpc.defaultie35CW.id
  zone_id    = var.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "defaultHrZmxC" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.defaultg0DCAR.id
  db_instance_storage = "20"
  vpc_id              = alicloud_vpc.defaultie35CW.id
  db_instance_class   = "mdb.shard.4x.large.d"
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = var.zone_id
}


`, name)
}

// Test Mongodb PrivateSrvNetworkAddress. <<< Resource test cases, automatically generated.
