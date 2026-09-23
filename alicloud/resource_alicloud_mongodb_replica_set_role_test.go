package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestCheckResourceMongoDBReplicaSetRoleExist should be prefixed by upper case, in other words, should be visible
// outside, otherwise it cannot be reflected by the test framework.
func (s *MongoDBService) TestCheckResourceMongoDBReplicaSetRoleExist(_ string) (interface{}, error) {
	return nil, nil
}

func TestAccAliCloudMongoDBReplicaSetRole_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_replica_set_role.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AlicloudMongoDBReplicaSetRoleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "TestCheckResourceMongoDBReplicaSetRoleExist")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBReplicaSetRole%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBReplicaSetRoleBasicDependence0)

	checkDestroy := func(state *terraform.State) error {
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
					"role_id":        "${alicloud_mongodb_instance.default.replica_sets[0].role_id}",
					"network_type":   "VPC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":    CHECKSET,
						"connection_prefix": CHECKSET,
						"connection_port":   CHECKSET,
						"replica_set_role":  CHECKSET,
						"connection_domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "test-mongodb-connection-modification-private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": CHECKSET,
						"connection_port":   CHECKSET,
						"replica_set_role":  CHECKSET,
						"connection_domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_port": "3729",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_port": "3729",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_port": "3739",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_port": "3739",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type":      "Public",
					"connection_prefix": "test-mongodb-connection-modification-public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": CHECKSET,
						"connection_port":   CHECKSET,
						"replica_set_role":  CHECKSET,
						"connection_domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_id":           "${alicloud_mongodb_public_network_address.default.replica_sets[1].role_id}",
					"connection_prefix": "test-mongodb-connection-modification-public-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": CHECKSET,
						"connection_port":   CHECKSET,
						"replica_set_role":  CHECKSET,
						"connection_domain": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_port": "3720",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_port": "3720",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "test-mongodb-connection-modification-public-2",
					"connection_port":   "3721",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "test-mongodb-connection-modification-public-2",
						"connection_port":   "3721",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
			},
		},
	})
}

var AlicloudMongoDBReplicaSetRoleMap0 = map[string]string{}

func AliCloudMongoDBReplicaSetRoleBasicDependence0(name string) string {
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

resource "alicloud_mongodb_public_network_address" "default" {
  db_instance_id = alicloud_mongodb_instance.default.id
}

`, name)
}
