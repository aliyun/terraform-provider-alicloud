package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestDescribeMongoDBPublicNetworkAddress should be visible, which can be reflected by the test framework.
func (s *MongoDBService) TestDescribeMongoDBPublicNetworkAddress(id string) (interface{}, error) {
	object, err := s.DescribeReplicaSetRole(id)
	if err != nil {
		return nil, WrapError(err)
	}

	if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
		if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
			for _, replicaSets := range replicaSetsList.([]interface{}) {
				replicaSetsArg := replicaSets.(map[string]interface{})
				networkType, ok := replicaSetsArg["NetworkType"]
				if ok && networkType == "Public" {
					return object, nil
				}
			}
		}
	}

	// simulate NotFoundError.
	return nil, &ComplexError{
		Err: fmt.Errorf("%s: public network address of %s", ResourceNotfound, id),
	}
}

func TestAccAliCloudMongoDBPublicNetworkAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_public_network_address.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AlicloudMongoDBPublicNetworkAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "TestDescribeMongoDBPublicNetworkAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBPublicNetworkAddress%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBPublicNetworkAddressBasicDependence0)

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
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"replica_sets.#": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudMongoDBPublicNetworkAddressMap0 = map[string]string{}

func AliCloudMongoDBPublicNetworkAddressBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_mongodb_zones" "zones_ids" {}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_mongodb_zones.zones_ids.ids[0]
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.default.id
  db_instance_storage = "20"
  vpc_id              = alicloud_vpc.default.id
  db_instance_class   = "mdb.shard.4x.large.d"
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = data.alicloud_mongodb_zones.zones_ids.ids[0]
}

`, name)
}
