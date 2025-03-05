package alicloud

import (
	"fmt"
	"log"
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

	var cachedReplicaSets map[string]map[string]interface{}
	var cacheReplicaSets resource.TestCheckFunc = func(state *terraform.State) error {
		id := state.RootModule().Resources[resourceId].Primary.ID
		ddsService := serverFunc().(*MongoDBService)

		object, err := ddsService.DescribeReplicaSetRole(id)
		if err != nil {
			return err
		}

		if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
			if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
				cachedReplicaSets = transferToReplicaSetMaps(replicaSetsList)
			}
		}
		log.Printf("[DEBUG] cachedReplicaSets: %v", cachedReplicaSets)
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
					cacheReplicaSets,
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

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone_id
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
  zone_id             = var.zone_id
}

`, name)
}
