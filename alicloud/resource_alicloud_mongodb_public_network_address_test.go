package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudMongoDBPublicNetworkAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_public_network_address.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AlicloudMongoDBPublicNetworkAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeReplicaSetRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBPublicNetworkAddress%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBPublicNetworkAddressBasicDependence0)

	checkDestroy := func(state *terraform.State) error {
		strs := strings.Split(rc.resourceId, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "alicloud_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}
		assert.NotEmpty(t, resourceType, "resourceType is empty")

		ddsService := serverFunc().(*MongoDBService)
		for _, rs := range state.RootModule().Resources {
			if rs.Type == resourceType {
				continue
			}

			object, err := ddsService.DescribeReplicaSetRole(rs.Primary.ID)
			if err != nil && NotFoundError(err) {
				continue
			} else if err != nil {
				return WrapError(err)
			}

			if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
				if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
					for _, replicaSets := range replicaSetsList.([]interface{}) {
						replicaSetsArg := replicaSets.(map[string]interface{})
						networkType, ok := replicaSetsArg["NetworkType"]
						if ok && networkType == "Public" {
							return WrapError(Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
						}
					}
				}
			}
		}

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
				ImportStateVerifyIgnore: []string{"db_instance_id"},
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

data "alicloud_mongodb_zones" "default" {
}

locals {
  index   = length(data.alicloud_mongodb_zones.default.zones) - 1
  zone_id = data.alicloud_mongodb_zones.default.zones[local.index].id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.zone_id
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  vswitch_id          = alicloud_vswitch.default.id
  name                = var.name
}
`, name)
}
