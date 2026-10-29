package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMongoDBPublicNetworkAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_public_network_address.default"
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMongoDBPublicNetworkAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "TestDescribeMongoDBPublicNetworkAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBPublicNetworkAddress%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBPublicNetworkAddressBasicDependence0)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": acctest.CHECKSET,
						"replica_sets.#": acctest.CHECKSET,
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
